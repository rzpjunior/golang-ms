# MODEL API EDENFARM VERSION 2

## Installation

1. pull `https://git.edenfarm.id/project-version2/datamodel.git`
2. place the repository in the same working directory as the primary repository like below:


```.
├── git.edenfarm.id
│   ├── project-version2
        ├── api
        ├── datamodel
 ```        

 ## How to use
 1. Modify the models as needed
 2. Navigate to your primary repository directory (i.e. api)
 3. Perform `go mod vendor` to refresh the latest changes to the `datamodel`
 4. Push your changes from both repositories `datamodel` and primary repository once the development is complete