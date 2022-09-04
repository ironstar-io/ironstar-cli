# gon.hcl
#
# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source = ["./dist/ironstar-cli_darwin_all/iron-macos"]
bundle_id = "io.ironstar.cli"

apple_id {
  username = "@env:AC_USERNAME"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Ironstar Hosting Services Pty Ltd (L7G23W3WF3)"
}
