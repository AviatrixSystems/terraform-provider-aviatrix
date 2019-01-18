""" Script to update vendor dependencies
    Needs: sudo -H pip install -U checksumdir  """
import os
import shutil
import json
import subprocess
import checksumdir


def get_immediate_subdirectories(a_dir):
    """ Get immediate subdirectories"""
    return [name for name in os.listdir(a_dir) if os.path.isdir(os.path.join(a_dir, name))]


def delete_everything_except(dir_name, exception_name, git=False):
    """Delete everything in folder_name except exception_name"""
    # Delete everything in $GOPATH/src except github.com
    print("Traversing" + dir_name)
    dir_list = get_immediate_subdirectories(dir_name)

    for dir_ in dir_list:
        if dir_ != exception_name:
            filename = os.path.join(dir_name, dir_)
            print("Deleting " + filename)
            if not git:
                shutil.rmtree(filename, ignore_errors=True)
            else:
                os.chdir(dir_name)
                subprocess.check_output("git rm -r " + filename, shell=True)


def move_everything_except(dir_name, exception_name, dest_dir, git=False):
    """Delete everything in folder_name except exception_name"""
    print("Traversing" + dir_name)
    dir_list = get_immediate_subdirectories(dir_name)

    for dir_ in dir_list:
        if dir_ not in exception_name:
            filename = os.path.join(dir_name, dir_)
            print("Moving " + filename)
            shutil.move(filename, dest_dir)
            if git:
                os.chdir(dest_dir)
                subprocess.check_output("git add " + dir_, shell=True)


def generate_vendor_json(github_path, vendor_json, origin_terr="",
                         site_name="github.com"):
    """ Generate vendor.json file"""
    dir_list = get_immediate_subdirectories(github_path)

    for dir_ in dir_list:
        for subdir_ in get_immediate_subdirectories(os.path.join(github_path, dir_)):
            git_folder = os.path.join(*[github_path, dir_, subdir_])
            os.chdir(git_folder)
            git_msg = subprocess.check_output("git log -1", shell=True)
            git_lines = git_msg.split("\n")
            commit = git_lines[0][7:]
            date_rev = git_lines[2][8:]
            check_sum = checksumdir.dirhash(git_folder, hashfunc='sha1')
            path = os.path.join(*[site_name, dir_, subdir_])
            print(path)
            origin = path if not origin_terr else origin_terr + path
            vendor_json["package"].append({"checksumSHA1": check_sum, "origin": origin,
                                           "path": path, "revision": commit,
                                           "revisionTime": date_rev})


def main():
    """Main function"""
    vendor_json = {"comment": "", "ignore": "test", "package": [],
                   "rootPath": "github.com/terraform-providers/terraform-provider-aviatrix"}
    gopath = os.environ["GOPATH"]
    terraform_path = os.path.join(*[gopath, "src", "github.com", "terraform-providers",
                                    "terraform-provider-aviatrix"])

    delete_everything_except(os.path.join(*[gopath, "src"]), exception_name="github.com")
    delete_everything_except(os.path.join(*[gopath, "src", "github.com"]),
                             exception_name="terraform-providers")
    delete_everything_except(os.path.join(*[terraform_path, "vendor"]), exception_name="github.com",
                             git=True)
    delete_everything_except(os.path.join(*[terraform_path, "vendor", "github.com"]),
                             exception_name="AviatrixSystems",
                             git=True)

    os.chdir(terraform_path)
    print("Obtaining latest dependancies using go get")
    os.system("go get")

    move_everything_except(os.path.join(*[gopath, "src"]), exception_name=["github.com"],
                           dest_dir=os.path.join(*[terraform_path, "vendor"]),
                           git=True)
    move_everything_except(os.path.join(*[gopath, "src", "github.com"]),
                           exception_name=["terraform-providers"],
                           dest_dir=os.path.join(*[terraform_path, "vendor", "github.com"]),
                           git=True)

    print("Dependencies\n-----------")
    generate_vendor_json(os.path.join(*[terraform_path, "vendor", "github.com"]), vendor_json)
    hashicorp_path = os.path.join(*[terraform_path, "vendor", "github.com", "hashicorp",
                                    "terraform", "vendor"])
    for dir_ in get_immediate_subdirectories(hashicorp_path):
        generate_vendor_json(os.path.join(hashicorp_path, dir_), vendor_json,
                             origin_terr="github.com/hashicorp/terraform/vendor/",
                             site_name=dir_)

    with open(os.path.join(*[terraform_path, "vendor", "vendor.json"]), 'w') as fileh:
        fileh.write(json.dumps(vendor_json, indent=2))
    print("Wrote vendor.json")

    os.chdir(terraform_path)
    move_everything_except(hashicorp_path, exception_name=["github.com"],
                           dest_dir=os.path.join(terraform_path, "vendor"), git=True)
    move_everything_except(os.path.join(hashicorp_path, "github.com"),
                           exception_name=["google", "hashicorp"],
                           dest_dir=os.path.join(*[terraform_path, "vendor", "github.com"]),
                           git=True)
    move_everything_except(os.path.join(hashicorp_path, "github.com", "google"),
                           exception_name=["go-querystring"],
                           dest_dir=os.path.join(*[terraform_path, "vendor", "github.com",
                                                   "google"]),
                           git=True)
    move_everything_except(os.path.join(hashicorp_path, "github.com", "hashicorp"),
                           exception_name=["****"],
                           dest_dir=os.path.join(*[terraform_path, "vendor", "github.com",
                                                   "hashicorp"]),
                           git=True)



    # os.system("git mv vendor/github.com/hashicorp/terraform/vendor/* vendor/")



if __name__ == "__main__":
    main()
