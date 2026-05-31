package collector

import (
	"testing"
)

func TestParseGpuNvidia(t *testing.T) {
	t.Parallel()

	data := []string{`<?xml version="1.0" ?>
<!DOCTYPE nvidia_smi_log SYSTEM "nvsmi_device_v12.dtd">
<nvidia_smi_log>
	<timestamp>Sun Feb 25 15:06:53 2024</timestamp>
	<driver_version>530.30.02</driver_version>
	<cuda_version>Not Found</cuda_version>
	<attached_gpus>1</attached_gpus>
	<gpu id="00000000:09:00.0">
		<product_name>NVIDIA GeForce GTX 1070</product_name>
		<product_brand>GeForce</product_brand>
		<product_architecture>Pascal</product_architecture>
		<display_mode>Disabled</display_mode>
		<display_active>Disabled</display_active>
		<persistence_mode>Disabled</persistence_mode>
		<mig_mode>
			<current_mig>N/A</current_mig>
			<pending_mig>N/A</pending_mig>
		</mig_mode>
		<mig_devices>
			None
		</mig_devices>
		<accounting_mode>Disabled</accounting_mode>
		<accounting_mode_buffer_size>4000</accounting_mode_buffer_size>
		<driver_model>
			<current_dm>N/A</current_dm>
			<pending_dm>N/A</pending_dm>
		</driver_model>
		<serial>N/A</serial>
		<uuid>GPU-0f28aa42-45b8-2019-1c8b-35a2ccd7ce89</uuid>
		<minor_number>0</minor_number>
		<vbios_version>86.04.50.00.24</vbios_version>
		<multigpu_board>No</multigpu_board>
		<board_id>0x900</board_id>
		<board_part_number>N/A</board_part_number>
		<gpu_part_number>1B81-200-A1</gpu_part_number>
		<gpu_fru_part_number>N/A</gpu_fru_part_number>
		<gpu_module_id>1</gpu_module_id>
		<inforom_version>
			<img_version>G001.0000.01.04</img_version>
			<oem_object>1.1</oem_object>
			<ecc_object>N/A</ecc_object>
			<pwr_object>N/A</pwr_object>
		</inforom_version>
		<gpu_operation_mode>
			<current_gom>N/A</current_gom>
			<pending_gom>N/A</pending_gom>
		</gpu_operation_mode>
		<gsp_firmware_version>N/A</gsp_firmware_version>
		<gpu_virtualization_mode>
			<virtualization_mode>None</virtualization_mode>
			<host_vgpu_mode>N/A</host_vgpu_mode>
		</gpu_virtualization_mode>
		<gpu_reset_status>
			<reset_required>No</reset_required>
			<drain_and_reset_recommended>N/A</drain_and_reset_recommended>
		</gpu_reset_status>
		<ibmnpu>
			<relaxed_ordering_mode>N/A</relaxed_ordering_mode>
		</ibmnpu>
		<pci>
			<pci_bus>09</pci_bus>
			<pci_device>00</pci_device>
			<pci_domain>0000</pci_domain>
			<pci_device_id>1B8110DE</pci_device_id>
			<pci_bus_id>00000000:09:00.0</pci_bus_id>
			<pci_sub_system_id>33011462</pci_sub_system_id>
			<pci_gpu_link_info>
				<pcie_gen>
					<max_link_gen>3</max_link_gen>
					<current_link_gen>3</current_link_gen>
					<device_current_link_gen>3</device_current_link_gen>
					<max_device_link_gen>3</max_device_link_gen>
					<max_host_link_gen>3</max_host_link_gen>
				</pcie_gen>
				<link_widths>
					<max_link_width>16x</max_link_width>
					<current_link_width>16x</current_link_width>
				</link_widths>
			</pci_gpu_link_info>
			<pci_bridge_chip>
				<bridge_chip_type>N/A</bridge_chip_type>
				<bridge_chip_fw>N/A</bridge_chip_fw>
			</pci_bridge_chip>
			<replay_counter>0</replay_counter>
			<replay_rollover_counter>0</replay_rollover_counter>
			<tx_util>0 KB/s</tx_util>
			<rx_util>4000 KB/s</rx_util>
			<atomic_caps_inbound>N/A</atomic_caps_inbound>
			<atomic_caps_outbound>N/A</atomic_caps_outbound>
		</pci>
		<fan_speed>49 %</fan_speed>
		<performance_state>P2</performance_state>
		<clocks_throttle_reasons>
			<clocks_throttle_reason_gpu_idle>Not Active</clocks_throttle_reason_gpu_idle>
			<clocks_throttle_reason_applications_clocks_setting>Not Active</clocks_throttle_reason_applications_clocks_setting>
			<clocks_throttle_reason_sw_power_cap>Not Active</clocks_throttle_reason_sw_power_cap>
			<clocks_throttle_reason_hw_slowdown>Not Active</clocks_throttle_reason_hw_slowdown>
			<clocks_throttle_reason_hw_thermal_slowdown>Not Active</clocks_throttle_reason_hw_thermal_slowdown>
			<clocks_throttle_reason_hw_power_brake_slowdown>Not Active</clocks_throttle_reason_hw_power_brake_slowdown>
			<clocks_throttle_reason_sync_boost>Not Active</clocks_throttle_reason_sync_boost>
			<clocks_throttle_reason_sw_thermal_slowdown>Not Active</clocks_throttle_reason_sw_thermal_slowdown>
			<clocks_throttle_reason_display_clocks_setting>Not Active</clocks_throttle_reason_display_clocks_setting>
		</clocks_throttle_reasons>
		<fb_memory_usage>
			<total>8192 MiB</total>
			<reserved>78 MiB</reserved>
			<used>1270 MiB</used>
			<free>6843 MiB</free>
		</fb_memory_usage>
		<bar1_memory_usage>
			<total>256 MiB</total>
			<used>2 MiB</used>
			<free>254 MiB</free>
		</bar1_memory_usage>
		<compute_mode>Default</compute_mode>
		<utilization>
			<gpu_util>2 %</gpu_util>
			<memory_util>1 %</memory_util>
			<encoder_util>3 %</encoder_util>
			<decoder_util>5 %</decoder_util>
		</utilization>
		<encoder_stats>
			<session_count>0</session_count>
			<average_fps>0</average_fps>
			<average_latency>0</average_latency>
		</encoder_stats>
		<fbc_stats>
			<session_count>0</session_count>
			<average_fps>0</average_fps>
			<average_latency>0</average_latency>
		</fbc_stats>
		<ecc_mode>
			<current_ecc>N/A</current_ecc>
			<pending_ecc>N/A</pending_ecc>
		</ecc_mode>
		<ecc_errors>
			<volatile>
				<single_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</single_bit>
				<double_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</double_bit>
			</volatile>
			<aggregate>
				<single_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</single_bit>
				<double_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</double_bit>
			</aggregate>
		</ecc_errors>
		<retired_pages>
			<multiple_single_bit_retirement>
				<retired_count>N/A</retired_count>
				<retired_pagelist>N/A</retired_pagelist>
			</multiple_single_bit_retirement>
			<double_bit_retirement>
				<retired_count>N/A</retired_count>
				<retired_pagelist>N/A</retired_pagelist>
			</double_bit_retirement>
			<pending_blacklist>N/A</pending_blacklist>
			<pending_retirement>N/A</pending_retirement>
		</retired_pages>
		<remapped_rows>N/A</remapped_rows>
		<temperature>
			<gpu_temp>60 C</gpu_temp>
			<gpu_temp_max_threshold>99 C</gpu_temp_max_threshold>
			<gpu_temp_slow_threshold>96 C</gpu_temp_slow_threshold>
			<gpu_temp_max_gpu_threshold>N/A</gpu_temp_max_gpu_threshold>
			<gpu_target_temperature>83 C</gpu_target_temperature>
			<memory_temp>N/A</memory_temp>
			<gpu_temp_max_mem_threshold>N/A</gpu_temp_max_mem_threshold>
		</temperature>
		<supported_gpu_target_temp>
			<gpu_target_temp_min>60 C</gpu_target_temp_min>
			<gpu_target_temp_max>92 C</gpu_target_temp_max>
		</supported_gpu_target_temp>
		<power_readings>
			<power_state>P2</power_state>
			<power_management>Supported</power_management>
			<power_draw>39.10 W</power_draw>
			<power_limit>190.00 W</power_limit>
			<default_power_limit>190.00 W</default_power_limit>
			<enforced_power_limit>190.00 W</enforced_power_limit>
			<min_power_limit>95.00 W</min_power_limit>
			<max_power_limit>200.00 W</max_power_limit>
		</power_readings>
		<clocks>
			<graphics_clock>1531 MHz</graphics_clock>
			<sm_clock>1531 MHz</sm_clock>
			<mem_clock>3802 MHz</mem_clock>
			<video_clock>1366 MHz</video_clock>
		</clocks>
		<applications_clocks>
			<graphics_clock>N/A</graphics_clock>
			<mem_clock>N/A</mem_clock>
		</applications_clocks>
		<default_applications_clocks>
			<graphics_clock>N/A</graphics_clock>
			<mem_clock>N/A</mem_clock>
		</default_applications_clocks>
		<deferred_clocks>
			<mem_clock>N/A</mem_clock>
		</deferred_clocks>
		<max_clocks>
			<graphics_clock>1936 MHz</graphics_clock>
			<sm_clock>1936 MHz</sm_clock>
			<mem_clock>4004 MHz</mem_clock>
			<video_clock>1708 MHz</video_clock>
		</max_clocks>
		<max_customer_boost_clocks>
			<graphics_clock>N/A</graphics_clock>
		</max_customer_boost_clocks>
		<clock_policy>
			<auto_boost>N/A</auto_boost>
			<auto_boost_default>N/A</auto_boost_default>
		</clock_policy>
		<voltage>
			<graphics_volt>N/A</graphics_volt>
		</voltage>
		<fabric>
			<state>N/A</state>
			<status>N/A</status>
		</fabric>
		<supported_clocks>
			<supported_mem_clock>
				<value>4004 MHz</value>
				<supported_graphics_clock>1936 MHz</supported_graphics_clock>
				<supported_graphics_clock>1923 MHz</supported_graphics_clock>
			</supported_mem_clock>
			<supported_mem_clock>
				<value>3802 MHz</value>
				<supported_graphics_clock>1936 MHz</supported_graphics_clock>
				<supported_graphics_clock>1923 MHz</supported_graphics_clock>
			</supported_mem_clock>
			<supported_mem_clock>
				<value>810 MHz</value>
				<supported_graphics_clock>1911 MHz</supported_graphics_clock>
				<supported_graphics_clock>1898 MHz</supported_graphics_clock>
			</supported_mem_clock>
			<supported_mem_clock>
				<value>405 MHz</value>
				<supported_graphics_clock>607 MHz</supported_graphics_clock>
				<supported_graphics_clock>594 MHz</supported_graphics_clock>
			</supported_mem_clock>
		</supported_clocks>
		<processes>
			<process_info>
				<gpu_instance_id>N/A</gpu_instance_id>
				<compute_instance_id>N/A</compute_instance_id>
				<pid>200870</pid>
				<type>C</type>
				<process_name></process_name>
				<used_memory>958 MiB</used_memory>
			</process_info>
			<process_info>
				<gpu_instance_id>N/A</gpu_instance_id>
				<compute_instance_id>N/A</compute_instance_id>
				<pid>1688347</pid>
				<type>C</type>
				<process_name></process_name>
				<used_memory>310 MiB</used_memory>
			</process_info>
		</processes>
		<accounted_processes>
		</accounted_processes>
	</gpu>
	<gpu id="00000000:a9:00.0">
		<product_name>NVIDIA GeForce GTX 3080 Ti</product_name>
		<product_brand>GeForce</product_brand>
		<product_architecture>Pascal</product_architecture>
		<display_mode>Disabled</display_mode>
		<display_active>Disabled</display_active>
		<persistence_mode>Disabled</persistence_mode>
		<mig_mode>
			<current_mig>N/A</current_mig>
			<pending_mig>N/A</pending_mig>
		</mig_mode>
		<mig_devices>
			None
		</mig_devices>
		<accounting_mode>Disabled</accounting_mode>
		<accounting_mode_buffer_size>4000</accounting_mode_buffer_size>
		<driver_model>
			<current_dm>N/A</current_dm>
			<pending_dm>N/A</pending_dm>
		</driver_model>
		<serial>N/A</serial>
		<uuid>GPU-0f28aa42-45b8-2019-1c8b-35a2ccd7ce89</uuid>
		<minor_number>0</minor_number>
		<vbios_version>86.04.50.00.24</vbios_version>
		<multigpu_board>No</multigpu_board>
		<board_id>0x900</board_id>
		<board_part_number>N/A</board_part_number>
		<gpu_part_number>1B81-200-A1</gpu_part_number>
		<gpu_fru_part_number>N/A</gpu_fru_part_number>
		<gpu_module_id>1</gpu_module_id>
		<inforom_version>
			<img_version>G001.0000.01.04</img_version>
			<oem_object>1.1</oem_object>
			<ecc_object>N/A</ecc_object>
			<pwr_object>N/A</pwr_object>
		</inforom_version>
		<gpu_operation_mode>
			<current_gom>N/A</current_gom>
			<pending_gom>N/A</pending_gom>
		</gpu_operation_mode>
		<gsp_firmware_version>N/A</gsp_firmware_version>
		<gpu_virtualization_mode>
			<virtualization_mode>None</virtualization_mode>
			<host_vgpu_mode>N/A</host_vgpu_mode>
		</gpu_virtualization_mode>
		<gpu_reset_status>
			<reset_required>No</reset_required>
			<drain_and_reset_recommended>N/A</drain_and_reset_recommended>
		</gpu_reset_status>
		<ibmnpu>
			<relaxed_ordering_mode>N/A</relaxed_ordering_mode>
		</ibmnpu>
		<pci>
			<pci_bus>12</pci_bus>
			<pci_device>01</pci_device>
			<pci_domain>0000</pci_domain>
			<pci_device_id>1B8110DE</pci_device_id>
			<pci_bus_id>00000000:a9:00.0</pci_bus_id>
			<pci_sub_system_id>33011462</pci_sub_system_id>
			<pci_gpu_link_info>
				<pcie_gen>
					<max_link_gen>3</max_link_gen>
					<current_link_gen>3</current_link_gen>
					<device_current_link_gen>3</device_current_link_gen>
					<max_device_link_gen>3</max_device_link_gen>
					<max_host_link_gen>3</max_host_link_gen>
				</pcie_gen>
				<link_widths>
					<max_link_width>16x</max_link_width>
					<current_link_width>16x</current_link_width>
				</link_widths>
			</pci_gpu_link_info>
			<pci_bridge_chip>
				<bridge_chip_type>N/A</bridge_chip_type>
				<bridge_chip_fw>N/A</bridge_chip_fw>
			</pci_bridge_chip>
			<replay_counter>0</replay_counter>
			<replay_rollover_counter>0</replay_rollover_counter>
			<tx_util>0 KB/s</tx_util>
			<rx_util>4000 KB/s</rx_util>
			<atomic_caps_inbound>N/A</atomic_caps_inbound>
			<atomic_caps_outbound>N/A</atomic_caps_outbound>
		</pci>
		<fan_speed>49 %</fan_speed>
		<performance_state>P2</performance_state>
		<clocks_throttle_reasons>
			<clocks_throttle_reason_gpu_idle>Not Active</clocks_throttle_reason_gpu_idle>
			<clocks_throttle_reason_applications_clocks_setting>Not Active</clocks_throttle_reason_applications_clocks_setting>
			<clocks_throttle_reason_sw_power_cap>Not Active</clocks_throttle_reason_sw_power_cap>
			<clocks_throttle_reason_hw_slowdown>Not Active</clocks_throttle_reason_hw_slowdown>
			<clocks_throttle_reason_hw_thermal_slowdown>Not Active</clocks_throttle_reason_hw_thermal_slowdown>
			<clocks_throttle_reason_hw_power_brake_slowdown>Not Active</clocks_throttle_reason_hw_power_brake_slowdown>
			<clocks_throttle_reason_sync_boost>Not Active</clocks_throttle_reason_sync_boost>
			<clocks_throttle_reason_sw_thermal_slowdown>Not Active</clocks_throttle_reason_sw_thermal_slowdown>
			<clocks_throttle_reason_display_clocks_setting>Not Active</clocks_throttle_reason_display_clocks_setting>
		</clocks_throttle_reasons>
		<fb_memory_usage>
			<total>8192 MiB</total>
			<reserved>78 MiB</reserved>
			<used>1270 MiB</used>
			<free>6843 MiB</free>
		</fb_memory_usage>
		<bar1_memory_usage>
			<total>256 MiB</total>
			<used>2 MiB</used>
			<free>254 MiB</free>
		</bar1_memory_usage>
		<compute_mode>Default</compute_mode>
		<utilization>
			<gpu_util>2 %</gpu_util>
			<memory_util>1 %</memory_util>
			<encoder_util>3 %</encoder_util>
			<decoder_util>5 %</decoder_util>
		</utilization>
		<encoder_stats>
			<session_count>0</session_count>
			<average_fps>0</average_fps>
			<average_latency>0</average_latency>
		</encoder_stats>
		<fbc_stats>
			<session_count>0</session_count>
			<average_fps>0</average_fps>
			<average_latency>0</average_latency>
		</fbc_stats>
		<ecc_mode>
			<current_ecc>N/A</current_ecc>
			<pending_ecc>N/A</pending_ecc>
		</ecc_mode>
		<ecc_errors>
			<volatile>
				<single_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</single_bit>
				<double_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</double_bit>
			</volatile>
			<aggregate>
				<single_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</single_bit>
				<double_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</double_bit>
			</aggregate>
		</ecc_errors>
		<retired_pages>
			<multiple_single_bit_retirement>
				<retired_count>N/A</retired_count>
				<retired_pagelist>N/A</retired_pagelist>
			</multiple_single_bit_retirement>
			<double_bit_retirement>
				<retired_count>N/A</retired_count>
				<retired_pagelist>N/A</retired_pagelist>
			</double_bit_retirement>
			<pending_blacklist>N/A</pending_blacklist>
			<pending_retirement>N/A</pending_retirement>
		</retired_pages>
		<remapped_rows>N/A</remapped_rows>
		<temperature>
			<gpu_temp>60 C</gpu_temp>
			<gpu_temp_max_threshold>99 C</gpu_temp_max_threshold>
			<gpu_temp_slow_threshold>96 C</gpu_temp_slow_threshold>
			<gpu_temp_max_gpu_threshold>N/A</gpu_temp_max_gpu_threshold>
			<gpu_target_temperature>83 C</gpu_target_temperature>
			<memory_temp>N/A</memory_temp>
			<gpu_temp_max_mem_threshold>N/A</gpu_temp_max_mem_threshold>
		</temperature>
		<supported_gpu_target_temp>
			<gpu_target_temp_min>60 C</gpu_target_temp_min>
			<gpu_target_temp_max>92 C</gpu_target_temp_max>
		</supported_gpu_target_temp>
		<power_readings>
			<power_state>P2</power_state>
			<power_management>Supported</power_management>
			<power_draw>39.10 W</power_draw>
			<power_limit>190.00 W</power_limit>
			<default_power_limit>190.00 W</default_power_limit>
			<enforced_power_limit>190.00 W</enforced_power_limit>
			<min_power_limit>95.00 W</min_power_limit>
			<max_power_limit>200.00 W</max_power_limit>
		</power_readings>
		<clocks>
			<graphics_clock>1531 MHz</graphics_clock>
			<sm_clock>1531 MHz</sm_clock>
			<mem_clock>3802 MHz</mem_clock>
			<video_clock>1366 MHz</video_clock>
		</clocks>
		<applications_clocks>
			<graphics_clock>N/A</graphics_clock>
			<mem_clock>N/A</mem_clock>
		</applications_clocks>
		<default_applications_clocks>
			<graphics_clock>N/A</graphics_clock>
			<mem_clock>N/A</mem_clock>
		</default_applications_clocks>
		<deferred_clocks>
			<mem_clock>N/A</mem_clock>
		</deferred_clocks>
		<max_clocks>
			<graphics_clock>1936 MHz</graphics_clock>
			<sm_clock>1936 MHz</sm_clock>
			<mem_clock>4004 MHz</mem_clock>
			<video_clock>1708 MHz</video_clock>
		</max_clocks>
		<max_customer_boost_clocks>
			<graphics_clock>N/A</graphics_clock>
		</max_customer_boost_clocks>
		<clock_policy>
			<auto_boost>N/A</auto_boost>
			<auto_boost_default>N/A</auto_boost_default>
		</clock_policy>
		<voltage>
			<graphics_volt>N/A</graphics_volt>
		</voltage>
		<fabric>
			<state>N/A</state>
			<status>N/A</status>
		</fabric>
		<processes>
			<process_info>
				<gpu_instance_id>N/A</gpu_instance_id>
				<compute_instance_id>N/A</compute_instance_id>
				<pid>200870</pid>
				<type>C</type>
				<process_name></process_name>
				<used_memory>958 MiB</used_memory>
			</process_info>
			<process_info>
				<gpu_instance_id>N/A</gpu_instance_id>
				<compute_instance_id>N/A</compute_instance_id>
				<pid>1688347</pid>
				<type>C</type>
				<process_name></process_name>
				<used_memory>310 MiB</used_memory>
			</process_info>
		</processes>
		<accounted_processes>
		</accounted_processes>
	</gpu>

</nvidia_smi_log>
`,
		`<?xml version="1.0" ?>
<!DOCTYPE nvidia_smi_log SYSTEM "nvsmi_device_v12.dtd">
<nvidia_smi_log>
	<timestamp>Wed Mar 13 13:06:02 2024</timestamp>
	<driver_version>550.54.14</driver_version>
	<cuda_version>Not Found</cuda_version>
	<attached_gpus>1</attached_gpus>
	<gpu id="00000000:09:00.0">
		<product_name>NVIDIA GeForce GTX 1070</product_name>
		<product_brand>GeForce</product_brand>
		<product_architecture>Pascal</product_architecture>
		<display_mode>Disabled</display_mode>
		<display_active>Disabled</display_active>
		<persistence_mode>Disabled</persistence_mode>
		<addressing_mode>N/A</addressing_mode>
		<mig_mode>
			<current_mig>N/A</current_mig>
			<pending_mig>N/A</pending_mig>
		</mig_mode>
		<mig_devices>
			None
		</mig_devices>
		<accounting_mode>Disabled</accounting_mode>
		<accounting_mode_buffer_size>4000</accounting_mode_buffer_size>
		<driver_model>
			<current_dm>N/A</current_dm>
			<pending_dm>N/A</pending_dm>
		</driver_model>
		<serial>N/A</serial>
		<uuid>GPU-0f28aa42-45b8-2019-1c8b-35a2ccd7ce89</uuid>
		<minor_number>0</minor_number>
		<vbios_version>86.04.50.00.24</vbios_version>
		<multigpu_board>No</multigpu_board>
		<board_id>0x900</board_id>
		<board_part_number>N/A</board_part_number>
		<gpu_part_number>1B81-200-A1</gpu_part_number>
		<gpu_fru_part_number>N/A</gpu_fru_part_number>
		<gpu_module_id>1</gpu_module_id>
		<inforom_version>
			<img_version>G001.0000.01.04</img_version>
			<oem_object>1.1</oem_object>
			<ecc_object>N/A</ecc_object>
			<pwr_object>N/A</pwr_object>
		</inforom_version>
		<inforom_bbx_flush>
			<latest_timestamp>N/A</latest_timestamp>
			<latest_duration>N/A</latest_duration>
		</inforom_bbx_flush>
		<gpu_operation_mode>
			<current_gom>N/A</current_gom>
			<pending_gom>N/A</pending_gom>
		</gpu_operation_mode>
		<c2c_mode>N/A</c2c_mode>
		<gpu_virtualization_mode>
			<virtualization_mode>None</virtualization_mode>
			<host_vgpu_mode>N/A</host_vgpu_mode>
			<vgpu_heterogeneous_mode>N/A</vgpu_heterogeneous_mode>
		</gpu_virtualization_mode>
		<gpu_reset_status>
			<reset_required>No</reset_required>
			<drain_and_reset_recommended>N/A</drain_and_reset_recommended>
		</gpu_reset_status>
		<gsp_firmware_version>N/A</gsp_firmware_version>
		<ibmnpu>
			<relaxed_ordering_mode>N/A</relaxed_ordering_mode>
		</ibmnpu>
		<pci>
			<pci_bus>09</pci_bus>
			<pci_device>00</pci_device>
			<pci_domain>0000</pci_domain>
			<pci_base_class>3</pci_base_class>
			<pci_sub_class>0</pci_sub_class>
			<pci_device_id>1B8110DE</pci_device_id>
			<pci_bus_id>00000000:09:00.0</pci_bus_id>
			<pci_sub_system_id>33011462</pci_sub_system_id>
			<pci_gpu_link_info>
				<pcie_gen>
					<max_link_gen>3</max_link_gen>
					<current_link_gen>1</current_link_gen>
					<device_current_link_gen>1</device_current_link_gen>
					<max_device_link_gen>3</max_device_link_gen>
					<max_host_link_gen>3</max_host_link_gen>
				</pcie_gen>
				<link_widths>
					<max_link_width>16x</max_link_width>
					<current_link_width>16x</current_link_width>
				</link_widths>
			</pci_gpu_link_info>
			<pci_bridge_chip>
				<bridge_chip_type>N/A</bridge_chip_type>
				<bridge_chip_fw>N/A</bridge_chip_fw>
			</pci_bridge_chip>
			<replay_counter>0</replay_counter>
			<replay_rollover_counter>0</replay_rollover_counter>
			<tx_util>0 KB/s</tx_util>
			<rx_util>0 KB/s</rx_util>
			<atomic_caps_inbound>N/A</atomic_caps_inbound>
			<atomic_caps_outbound>N/A</atomic_caps_outbound>
		</pci>
		<fan_speed>33 %</fan_speed>
		<performance_state>P8</performance_state>
		<clocks_event_reasons>
			<clocks_event_reason_gpu_idle>Active</clocks_event_reason_gpu_idle>
			<clocks_event_reason_applications_clocks_setting>Not Active</clocks_event_reason_applications_clocks_setting>
			<clocks_event_reason_sw_power_cap>Not Active</clocks_event_reason_sw_power_cap>
			<clocks_event_reason_hw_slowdown>Not Active</clocks_event_reason_hw_slowdown>
			<clocks_event_reason_hw_thermal_slowdown>Not Active</clocks_event_reason_hw_thermal_slowdown>
			<clocks_event_reason_hw_power_brake_slowdown>Not Active</clocks_event_reason_hw_power_brake_slowdown>
			<clocks_event_reason_sync_boost>Not Active</clocks_event_reason_sync_boost>
			<clocks_event_reason_sw_thermal_slowdown>Not Active</clocks_event_reason_sw_thermal_slowdown>
			<clocks_event_reason_display_clocks_setting>Not Active</clocks_event_reason_display_clocks_setting>
		</clocks_event_reasons>
		<sparse_operation_mode>N/A</sparse_operation_mode>
		<fb_memory_usage>
			<total>8192 MiB</total>
			<reserved>85 MiB</reserved>
			<used>984 MiB</used>
			<free>7122 MiB</free>
		</fb_memory_usage>
		<bar1_memory_usage>
			<total>256 MiB</total>
			<used>2 MiB</used>
			<free>254 MiB</free>
		</bar1_memory_usage>
		<cc_protected_memory_usage>
			<total>0 MiB</total>
			<used>0 MiB</used>
			<free>0 MiB</free>
		</cc_protected_memory_usage>
		<compute_mode>Default</compute_mode>
		<utilization>
			<gpu_util>0 %</gpu_util>
			<memory_util>0 %</memory_util>
			<encoder_util>0 %</encoder_util>
			<decoder_util>0 %</decoder_util>
			<jpeg_util>N/A</jpeg_util>
			<ofa_util>N/A</ofa_util>
		</utilization>
		<encoder_stats>
			<session_count>0</session_count>
			<average_fps>0</average_fps>
			<average_latency>0</average_latency>
		</encoder_stats>
		<fbc_stats>
			<session_count>0</session_count>
			<average_fps>0</average_fps>
			<average_latency>0</average_latency>
		</fbc_stats>
		<ecc_mode>
			<current_ecc>N/A</current_ecc>
			<pending_ecc>N/A</pending_ecc>
		</ecc_mode>
		<ecc_errors>
			<volatile>
				<single_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</single_bit>
				<double_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</double_bit>
			</volatile>
			<aggregate>
				<single_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</single_bit>
				<double_bit>
					<device_memory>N/A</device_memory>
					<register_file>N/A</register_file>
					<l1_cache>N/A</l1_cache>
					<l2_cache>N/A</l2_cache>
					<texture_memory>N/A</texture_memory>
					<texture_shm>N/A</texture_shm>
					<cbu>N/A</cbu>
					<total>N/A</total>
				</double_bit>
			</aggregate>
		</ecc_errors>
		<retired_pages>
			<multiple_single_bit_retirement>
				<retired_count>N/A</retired_count>
				<retired_pagelist>N/A</retired_pagelist>
			</multiple_single_bit_retirement>
			<double_bit_retirement>
				<retired_count>N/A</retired_count>
				<retired_pagelist>N/A</retired_pagelist>
			</double_bit_retirement>
			<pending_blacklist>N/A</pending_blacklist>
			<pending_retirement>N/A</pending_retirement>
		</retired_pages>
		<remapped_rows>N/A</remapped_rows>
		<temperature>
			<gpu_temp>39 C</gpu_temp>
			<gpu_temp_tlimit>N/A</gpu_temp_tlimit>
			<gpu_temp_max_threshold>99 C</gpu_temp_max_threshold>
			<gpu_temp_slow_threshold>96 C</gpu_temp_slow_threshold>
			<gpu_temp_max_gpu_threshold>N/A</gpu_temp_max_gpu_threshold>
			<gpu_target_temperature>83 C</gpu_target_temperature>
			<memory_temp>N/A</memory_temp>
			<gpu_temp_max_mem_threshold>N/A</gpu_temp_max_mem_threshold>
		</temperature>
		<supported_gpu_target_temp>
			<gpu_target_temp_min>60 C</gpu_target_temp_min>
			<gpu_target_temp_max>92 C</gpu_target_temp_max>
		</supported_gpu_target_temp>
		<gpu_power_readings>
			<power_state>P8</power_state>
			<power_draw>8.19 W</power_draw>
			<current_power_limit>190.00 W</current_power_limit>
			<requested_power_limit>190.00 W</requested_power_limit>
			<default_power_limit>190.00 W</default_power_limit>
			<min_power_limit>95.00 W</min_power_limit>
			<max_power_limit>200.00 W</max_power_limit>
		</gpu_power_readings>
		<gpu_memory_power_readings>
			<power_draw>N/A</power_draw>
		</gpu_memory_power_readings>
		<module_power_readings>
			<power_state>P8</power_state>
			<power_draw>N/A</power_draw>
			<current_power_limit>N/A</current_power_limit>
			<requested_power_limit>N/A</requested_power_limit>
			<default_power_limit>N/A</default_power_limit>
			<min_power_limit>N/A</min_power_limit>
			<max_power_limit>N/A</max_power_limit>
		</module_power_readings>
		<clocks>
			<graphics_clock>139 MHz</graphics_clock>
			<sm_clock>139 MHz</sm_clock>
			<mem_clock>405 MHz</mem_clock>
			<video_clock>544 MHz</video_clock>
		</clocks>
		<applications_clocks>
			<graphics_clock>N/A</graphics_clock>
			<mem_clock>N/A</mem_clock>
		</applications_clocks>
		<default_applications_clocks>
			<graphics_clock>N/A</graphics_clock>
			<mem_clock>N/A</mem_clock>
		</default_applications_clocks>
		<deferred_clocks>
			<mem_clock>N/A</mem_clock>
		</deferred_clocks>
		<max_clocks>
			<graphics_clock>1936 MHz</graphics_clock>
			<sm_clock>1936 MHz</sm_clock>
			<mem_clock>4004 MHz</mem_clock>
			<video_clock>1708 MHz</video_clock>
		</max_clocks>
		<max_customer_boost_clocks>
			<graphics_clock>N/A</graphics_clock>
		</max_customer_boost_clocks>
		<clock_policy>
			<auto_boost>N/A</auto_boost>
			<auto_boost_default>N/A</auto_boost_default>
		</clock_policy>
		<voltage>
			<graphics_volt>N/A</graphics_volt>
		</voltage>
		<fabric>
			<state>N/A</state>
			<status>N/A</status>
			<cliqueId>N/A</cliqueId>
			<clusterUuid>N/A</clusterUuid>
			<health>
				<bandwidth>N/A</bandwidth>
			</health>
		</fabric>
		<supported_clocks>
			<supported_mem_clock>
				<value>4004 MHz</value>
				<supported_graphics_clock>1936 MHz</supported_graphics_clock>
				<supported_graphics_clock>1923 MHz</supported_graphics_clock>
			</supported_mem_clock>
			<supported_mem_clock>
				<value>3802 MHz</value>
				<supported_graphics_clock>1936 MHz</supported_graphics_clock>
				<supported_graphics_clock>1923 MHz</supported_graphics_clock>
			</supported_mem_clock>
			<supported_mem_clock>
				<value>810 MHz</value>
				<supported_graphics_clock>1911 MHz</supported_graphics_clock>
				<supported_graphics_clock>1898 MHz</supported_graphics_clock>
			</supported_mem_clock>
			<supported_mem_clock>
				<value>405 MHz</value>
				<supported_graphics_clock>607 MHz</supported_graphics_clock>
				<supported_graphics_clock>594 MHz</supported_graphics_clock>
			</supported_mem_clock>
		</supported_clocks>
		<processes>
			<process_info>
				<gpu_instance_id>N/A</gpu_instance_id>
				<compute_instance_id>N/A</compute_instance_id>
				<pid>281698</pid>
				<type>C</type>
				<process_name></process_name>
				<used_memory>982 MiB</used_memory>
			</process_info>
		</processes>
		<accounted_processes>
		</accounted_processes>
	</gpu>

</nvidia_smi_log>
`}

	for k, d := range data {
		g := (&Collector{}).parseGpuNvidia(d)
		if k == 0 {
			if len(g) != 2 {
				t.Fatal("expected 2 GPUs, arr test data:", k)
			}
		}
		for _, gpu := range g {
			if k == 0 {
				if gpu.Power != "39.10 W" {
					t.Fatal("expected '39.10 W', arr test data:", k)
				}
				if gpu.PowerLimit != "190.00 W" {
					t.Fatal("expected '190.00 W', arr test data:", k)
				}
			}

			if k == 1 {
				if gpu.Power != "8.19 W" {
					t.Fatal("expected '8.19 W', arr test data:", k)
				}
				if gpu.PowerLimit != "190.00 W" {
					t.Fatal("expected '190.00 W', arr test data:", k)
				}
			}
		}
	}
}
