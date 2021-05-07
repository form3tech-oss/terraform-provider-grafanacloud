# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

A full example can be found in directory **full/**.

The document generation tool (`tfplugindocs`) looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. Other examples are included (i.e `./full/`) that serve as additional documentation, although not used by the documentation tool.

* **provider/provider.tf** example file for the provider index page
* **data-sources/\<full data source name\>/data-source.tf** example file for the named data source page
* **resources/\<full resource name\>/resource.tf** example file for the named data source page
