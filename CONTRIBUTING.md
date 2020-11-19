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
