## Contributing to the Aviatrix Terraform Provider
This is a collection of useful knowledge for making great contributions to the
provider.

### Provider Development Checklist
Before opening a new PR, go through each of the bullet points and ask yourself
if you have thought through all the steps that need to happen to support your
code change.
- Boilerplate
	- Acceptance Tests: Is your new feature/resource/attribute covered by an acceptance test?
	- Documentation: Have you updated the relevant doc page?
	- HCL Formatting: Is the HCL in your doc examples and acceptance tests formatted properly?
	- gofmt & goimports: Have you ran `make fmt imports` yet?
- Manual testing
	- Create: Verified possible creation configurations?
	- Update: Verified possible update configurations?
	- Delete: Verified deletion with different configurations?
	- Version upgrade: If modifying/adding an attribute, can an existing user still use their resource after upgrading?
	- Import: If you remove a resource from state, are you able to import it completely back into state?
- Other considerations
	- Attributes should be ordered in the schema definition like so:
		- Required
		- Optional
		- Optional and Computed
		- Computed 
	- In the Read function, if the resource does not exist do not return error, instead do the following (documented here https://learn.hashicorp.com/tutorials/terraform/provider-setup):
	```
	if err == goaviatrix.ErrNotFound {
		d.SetId("")
		return nil
	}
	```
	- If the Create or Update function performs more than one API call, the Read function needs to be called before returning. For example, `resourceAviatrixDeviceTagUpdate` can call up to 3 different APIs: `UpdateDeviceTagConfig`, `AttachDeviceTag` and `CommitDeviceTag`. The easiest way to ensure the Read is called is to use `defer`, in `resourceAviatrixDeviceTagUpdate` this is the first line of the function.
	```
	func resourceAviatrixDeviceTagUpdate(d *schema.ResourceData, meta interface{}) error {
		defer resourceAviatrixDeviceTagRead(d, meta)
	```
