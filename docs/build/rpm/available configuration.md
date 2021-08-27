# Package structure: Available options for config.yaml
___
# Meta:
| Name  | Description  | Type  | Yaml Variable |
|-------|--------------|-------|---------------|
| Name  | Package name |  string |name   |   |
| Version  | Package version   | string  | version  |   |
| Release  | OS release  | string  | release  |   |
| Architecture  | CPU Architecture  | string  | arch  |   |
| Summary  | Short description   | string  | summary  |   |
| Description  | Long description  | string  | description  |   |
| OS  | OS name  | string  | os  |   |
| Vendor  |  Vendor name | string  | vendor  |   |
| Url | Maintainer URL  | string  | url  |   |
| License  | Package License  | string  | license  |   |
| Maintainer  | Maintainer Information  | string  | maintainer  |   |
| Signature  |Package signature   |   | signature  |   |
| PrivateKey  |Sub variable of Signature   | string  | private_key  |   |
| PassPhrase  |Sub variable of Signature   | string  | pass_phrase  |   |

# File structure
| Name  | Description  | Type  | Yaml Variable |
|-------|--------------|-------|---------------|
| Destination  | Target destination where file will be placed |  string |destination   |   |
| Source  | Source file   | string  | source  |   |
| Body  | Immediate script or any data written inline  | string  | body  |   |
| Mode  | File mode, requires for Body, source will take mode directly from source file  | string  | mode  |   |
| Owner  | File owner   | string  | owner  |   |
| Group  | File group  | string  | group  |   |
| MTime  | File modification time, if empty it takes from file, in case body - takes current time. example: 2021-08-18 21:30:00 | time  | mtime  |   |

# File options
| Name     | Type      | Yaml Variable |
|----------|-----------|---------------|
| Generic  | map[string][]File structure | generic   |   |
| Config   | map[string][]File structure  | config  |   |
| Not_use  | map[string][]File structure  | not_use  |   |
| Missing_ok  | map[string][]File structure  | missing_ok  |   |
| No_replace  | map[string][]File structure  | no_replace  |   |
| Spec  | map[string][]File structure  | spec  |   |
| Ghost  | map[string][]File structure  | ghost  |   |
| License  | map[string][]File structure  | license  |   |
| Readme  | map[string][]File structure  | readme  |   |
| Exclude  | map[string][]File structure  | exclude  |   |

# Directory structure
| Name  | Description  | Type  | Yaml Variable |
|-------|--------------|-------|---------------|
| Destination  | Target destination directory, will be created |  string |destination   |   |
| Mode  | Directory mode | string  | mode  |   |
| Owner  | Directory owner   | string  | owner  |   |
| Group  | Directory group  | string  | group  |   |

# Directory options
| Name     | Type      | Yaml Variable |
|----------|-----------|---------------|
| Directory  | []Directory structure | directory   |   |

# PreInstall and PostInstall
| Name  | Description  | Type  | Yaml Variable |
|-------|--------------|-------|---------------|
| PreInstall  | PreInstallation activity, like create user etc. |  []string |preinstall   |   |
| PostInstall  | PostInstallation activity required for package  | []string  |postinstall  |   |
| PreUninstall  | PreUninstall activity required to remove package   | []string  |preuninstall |   |
| PostUninstall  | PostUninstall activity required to remove package  | []string  |postuninstall  |   |

# Dependency structure
| Name  | Description  | Type  | Yaml Variable |
|-------|--------------|-------|---------------|
| Name  | Component name, example kernel  |  string |name   |   |
| Version  | kernel or library version |  string |version   |   |
| Operator  | comparing operator > < >= <= |  string |version   |   |

# Dependency option
| Name  | Description  | Type  | Yaml Variable |
|-------|--------------|-------|---------------|
| Dependencies  | Dependencies definition  |  []Dependency structure |dependencies   |   |