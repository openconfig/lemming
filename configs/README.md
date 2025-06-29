# Lemming Device Configuration

Lemming can simulate various network device behaviors. This guide explains how to configure it.

## How to Configure Lemming

You have four options, listed in order of precedence:

**1. Command-Line Flag**

Specify a full path to your configuration file.

```bash
./lemming --config_file /path/to/your_config.textproto
```

**2. Environment Variable (`LEMMING_CONFIG_FILE`)**

This is the most flexible option.

* **Use a Vendor Preset:** Set the variable to a known vendor short name to use a pre-packaged configuration.

    ```bash
    # Use the built-in Arista configuration
    export LEMMING_CONFIG_FILE=arista
    ./lemming

    # Supported presets: arista, cisco, juniper
    ```

* **Use a File Path:** Set it to a full or relative path to your own configuration file.

    ```bash
    export LEMMING_CONFIG_FILE=./my_custom_config.textproto
    ./lemming
    ```

**3. Default File**

Place a file named `lemming_default.textproto` in the `configs/` directory. Lemming will automatically load it if no other configuration is specified.

**4. No Configuration**

If none of the above are provided, Lemming will start with its built-in default values.

## Configuration Overview

Configuration is done using `.textproto` files. You only need to specify the values you want to override; the rest will use defaults.

### Example: Custom Vendor ID

Create a file `my_config.textproto` with only the vendor info:

```protobuf
# my_config.textproto
vendor {
  name: "MyCustomDevice"
  model: "VirtualRouter-X1"
  os_version: "2.0-beta"
}
```

Then run Lemming with it:

```bash
./lemming --config_file my_config.textproto
```

All other settings (components, timing, etc.) will use the default values.

### Key Configuration Sections

You can customize the following parts of the device:

* **`vendor`**: The device's identity (e.g., name, model, OS version).
* **`components`**: The physical layout (e.g., number and names of line cards, supervisors).
* **`processes`**: Mock system processes to simulate for monitoring.
* **`timing`**: Durations for operations like reboots and switchovers.
* **`network_sim`**: Network behavior for ping, including latency and packet loss.

For detailed structure, see the `lemming_default.textproto` file.
