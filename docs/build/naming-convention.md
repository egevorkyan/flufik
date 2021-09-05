# RPM and DEB naming convention

><b>Note:</b> It is very important to know naming convention to build proper versioning for your future packages, to do
> not have issues with updates.

- RPM naming convention (example: flufik-0.1.1.x86_64.rpm) consists:
  1. Package name
  2. Package version
  3. Package release
  4. Package Architecture - for what architecture package was built
  5. Package extension

- DEB naming convention (example: flufik_0.1-1_amd64.deb) consists:
  1. Package name
  2. Package version
  3. Package release
  4. Package Architecture - for what architecture package was built
  5. Package extension

# How repositories considering if package has updated version to propose end host?
As an example: flufik-0.1.1.x86_64.rpm or flufik_0.1-1_amd64.deb

1. If new package version is higher than existing one
2. If new package version has same version, than release will be compared, as an example:
   - flufik's package version is 0.1
   - Current release is 1
   - If I add some fixes or very small changes and will build flufik with release 2, like flufik-0.1.2.x86_64.rpm,
     than flufik-0.1.2.x86_64.rpm will be considered as flufik's package update and repository server will propagate update 
     to end host.
3. <b>Warning:</b> If we don't have release and only have version which is violation of official naming convention than update will just replace existing package.