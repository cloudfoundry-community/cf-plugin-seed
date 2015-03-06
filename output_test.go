package main

const appSummary = `{
   "guid": "2dc79bb7-752e-47ef-bcd9-b3a2320a8474",
   "name": "space3",
   "apps": [
      {
         "guid": "6b3276d0-2697-4b89-b845-5bf942d8f06f",
         "urls": [
            "app3.10.244.0.34.xip.io"
         ],
         "routes": [
            {
               "guid": "25fea092-3155-47c9-a776-d6f6e688ae9c",
               "host": "app3",
               "domain": {
                  "guid": "f9ca5dc1-9c04-4da0-bb04-6e0849be3ddf",
                  "name": "10.244.0.34.xip.io"
               }
            }
         ],
         "service_count": 0,
         "service_names": [],
         "running_instances": 0,
         "name": "app3",
         "production": false,
         "space_guid": "2dc79bb7-752e-47ef-bcd9-b3a2320a8474",
         "stack_guid": "4130fde8-3544-4f8a-8306-760b1220ef14",
         "buildpack": null,
         "detected_buildpack": null,
         "environment_json": {},
         "memory": 256,
         "instances": 1,
         "disk_quota": 1024,
         "state": "STOPPED",
         "version": "7f2ba892-abbb-4536-a9b4-1143c19b0e58",
         "command": null,
         "console": false,
         "debug": null,
         "staging_task_id": null,
         "package_state": "PENDING",
         "health_check_type": "port",
         "health_check_timeout": null,
         "staging_failed_reason": null,
         "docker_image": null,
         "package_updated_at": null,
         "detected_start_command": ""
      },
      {
         "guid": "59de10d6-c292-4e53-9fc3-8e7047cdc201",
         "urls": [
            "app4.10.244.0.34.xip.io"
         ],
         "routes": [
            {
               "guid": "34e8328b-4434-44bf-8672-383149a263bc",
               "host": "app4",
               "domain": {
                  "guid": "f9ca5dc1-9c04-4da0-bb04-6e0849be3ddf",
                  "name": "10.244.0.34.xip.io"
               }
            }
         ],
         "service_count": 0,
         "service_names": [],
         "running_instances": 0,
         "name": "app4",
         "production": false,
         "space_guid": "2dc79bb7-752e-47ef-bcd9-b3a2320a8474",
         "stack_guid": "4130fde8-3544-4f8a-8306-760b1220ef14",
         "buildpack": null,
         "detected_buildpack": null,
         "environment_json": {},
         "memory": 256,
         "instances": 1,
         "disk_quota": 1024,
         "state": "STOPPED",
         "version": "71b8f549-53dd-4928-8987-cea7b86f863f",
         "command": null,
         "console": false,
         "debug": null,
         "staging_task_id": null,
         "package_state": "PENDING",
         "health_check_type": "port",
         "health_check_timeout": null,
         "staging_failed_reason": null,
         "docker_image": null,
         "package_updated_at": null,
         "detected_start_command": ""
      }
   ],
   "services": []
}`

const appSummaryNoRoute = `{
   "guid": "2dc79bb7-752e-47ef-bcd9-b3a2320a8474",
   "name": "space3",
   "apps": [
      {
         "guid": "6b3276d0-2697-4b89-b845-5bf942d8f06f",
         "urls": [
            "app3.10.244.0.34.xip.io"
         ],
         "routes": [],
         "service_count": 0,
         "service_names": [],
         "running_instances": 0,
         "name": "app3",
         "production": false,
         "space_guid": "2dc79bb7-752e-47ef-bcd9-b3a2320a8474",
         "stack_guid": "4130fde8-3544-4f8a-8306-760b1220ef14",
         "buildpack": null,
         "detected_buildpack": null,
         "environment_json": {},
         "memory": 256,
         "instances": 1,
         "disk_quota": 1024,
         "state": "STOPPED",
         "version": "7f2ba892-abbb-4536-a9b4-1143c19b0e58",
         "command": null,
         "console": false,
         "debug": null,
         "staging_task_id": null,
         "package_state": "PENDING",
         "health_check_type": "port",
         "health_check_timeout": null,
         "staging_failed_reason": null,
         "docker_image": null,
         "package_updated_at": null,
         "detected_start_command": ""
      }
   ],
   "services": []
}`
