# Generator

This software parse `doc` bluez folder and output a set of struct to interact with the bluez DBus API.


## Notes

- Generated files have a `gen_` prefix, followed by the API name
- If a `<API name>.go` file exists, it will be skipped from the generation. This to allow custom code to live with generated one.
- Generation process does not overwrite existing files, ensure to remove previously generated files.
