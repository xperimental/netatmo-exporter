# Introduction

# System setup

The docker configuration works great locally on my MacBook Pro as well as on a passive Intel based system under Ubuntu 20.4. The Ubuntu system is available in my local network (not externally). A tablet in kiosk mode is displaying the Grafana dashboard. The *docker-compose.yml* as well as all other configuration files are available in the *data* directory.  

# The Docker containers

The *docker-compose* file provided here loads 4 containers to visualize the data of your netatmo devices.  

## 1. The netatmo-exporter

Available via [Github](https://github.com/xperimental/netatmo-exporter). The container will connect to the Netatmo webserver reading the data of your devices available throughout your private NETATMO_CLIENT_ID ( updating every 8 minutes). The data will be available via a JSON file exposed by a webservice.  

## 2. A Prometheus data monitoring

Available via [Github](https://grafana.com/docs/grafana/latest/getting-started/get-started-grafana-prometheus/). Prometheus can monitor data provided by a webservice such as the *netatmo-exporter*. The data being captured can be visualized in a *Grafana* dashboard.

## 3. Grafana

Available via [Github](https://grafana.com). Grafana is being used to visualize any kind of data of different sources such as Prometheus. Within Grafana you create dashboards containing different styles of panels showing different type of data.  

## 4. Watchtower
A tool to monitor the version of your docker container to update them accordingly.  

The different container talk to each other via webservices and so they require a common network which is defined in the *docker-compose.yml* file. Do not use the IP addresses as part of your home network, instead define a separate network.  

# The docker-compose.yml

Prometheus as well as the netatmo-exporter and Grafana need different, specific configuration files. To make it easier to migrate the project to another host we bind a local directory to all of the containers. In this case the containers will load their configuration from one directory.
The file *.env* will be used by docker to define environment variables that are available throughout all containers defined in *docker-compose.yml*. The file has to exist in the same directory as the *docker-compose.yml* file.  

This is the content of *.env* file used here, fell free to modify the mapped directory according to your need:  

```shell
DATA_DIRECTORY="~/Documents/GitHub/netatmo-exporter/dockerExample/data"
```

Grafana stores its configuration in *grafana.db*. It turned out that using the same configuration file via a network share to compose the same project on different machines ended up in problems. As a workaround we create a docker volume for the Grafana configuration. Having it in a docker volume enables us to download the file via docker for backup purposes.  To create the docker volume we specify it on the top level of *docker-compose.yml*.  

```yml
volumes:
  grafana_data:
    driver: local
```

As said above we need to define our own network within the *docker-compose.yml* file. As an IP range we use 198.2.3.0/24 and define the network on the top-level of the file:  

```yml
networks:
  netatmo-network:
    driver: bridge
    ipam:
      config:
        - subnet: 192.168.3.0/24
  bridge:
    driver: bridge
```

Now we defined our network and a data directory for the Grafana configuration and can continue defining our containers.

## Prometheus configuration

To configure Prometheus we place the Prometheus configuration file called *prometheus.yml* in the data directory we specified in the *.env* file. This is a default *prometheus.yml* file:

```yml
global:
  scrape_interval: 5s 
  scrape_timeout: 5s
  external_labels:
    monitor: 'codelab-monitor'

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
    - targets: ['localhost:9090']

  - job_name: 'netatmo'
    scrape_interval: 45s
    static_configs:
      - targets: ['192.168.3.52:9210']
```

Notice that we use 192.168.3.52:9210 for the IP address and port for the *netatmo-exporter* instance to collect the data from. That of course has to be the same for the *netatmo_exporter* configuration in the *docker-compose.yml* file. Prometheus itself is using port 9090 and that needs to match our Prometheus instance defined in *docker-compose.yml*.

## Prometheus docker-compose configuration

All systems are defined under the *services* section in the *docker-compose.yml* file.

```yml
services:
  prometheus:
    # prometheus is scraping the netatmo exporter service in a given intervall. the configuration
    # is done by prometheus.yml. the file in this example is given in the ${DATA_DIRECTORY} directory
    image: prom/prometheus:latest
    volumes:
      # mount the data directory to /etc/prometheus to load the prometheus config from prometheus.yml
      - ${DATA_DIRECTORY}:/etc/prometheus    
    ports:
      # the default port being used by prometheus to send the data
      - "9090:9090"
    networks:
      netatmo-network:
        ipv4_address: 192.168.3.50
```

 We make use of the latest Prometheus docker package and we map the *DATA_DIRECTORY* to */etc/prometheus* within the container. Therefore Prometheus will use our configuration file as shown above. We define the port 9090 and assign the IP 192.168.3.50 of our network as defined in the *networks* section (refer above) to it.  

## Grafana docker-compose configuration  

This is straight forward. I added both methods for the configuration file and I assume both of them are working. But as I was playing with different configurations and systems it turned out that only the variant with a docker volume worked. Again the Grafana settings are below the *services* section.  

```yml
grafana:
    # grafana is running a dashboard collecting data from prometheus 
    image: grafana/grafana:latest
    volumes:
      # create a grafana directory that is stored locally within all images defined in this file
      - grafana_data:/var/lib/grafana
      # ...or use a mounted directory instead
      #- ${DATA_DIRECTORY}:/var/lib/grafana
    ports:
      # the default port beeing used by grafana to show the dashboard and to configure it
      - "3000:3000"
    networks:
      netatmo-network:
        ipv4_address: 192.168.3.51
      bridge:
```

We again make use of the latest Grafana docker package and this time we map the docker volume to the */var/lib/grafana* directory within the container. That's the directory where Grafana expects its configuration file. You can try to map the *DATA_DIRECTORY* to */var/lib/grafana* and I expect that it will work as well. We define the port 3000 and assign the IP 192.168.3.51 of our network as defined in the networks (refer above) to it.  

## Netatmo-exporter docker-compose configuration

Next is the *netatmo-exporter* container. The configuration is a bit longer as we also need to provide it with our netatmo credentials. Again the netatmo-exporter settings are below the *services* section.  

```yml
netatmo-exporter:
    # this image is loading data from the netatmo server and provides them as a webservice that is scraped by prometheus
    image: ghcr.io/xperimental/netatmo-exporter:latest
    restart: unless-stopped
    environment:
      - NETATMO_CLIENT_ID=your_client_id
      - NETATMO_CLIENT_SECRET=your_client_secret
      - NETATMO_EXPORTER_EXTERNAL_URL=http://192.168.3.52:9210
      - NETATMO_EXPORTER_TOKEN_FILE=/usr/token.json
    volumes:
      # the netatmo-exporter need your netatmo credentials stored in a filename named token.json. the file includes
      # your access_token, refresh_token and the expiry date that is refreshed automatically
      - ${DATA_DIRECTORY}:/usr
    ports:
        # the default port being used by prometheus to scrape the data
      - "9210:9210"
    networks:
      netatmo-network:
        ipv4_address: 192.168.3.52
      bridge:
```

We make again use of the latest *netatmo-exporter* docker package and we map the *DATA_DIRECTORY* to */usr* and so *netatmo-exporter* loads *token.json* from our *DATA_DIRECTORY*. Please apply your personal Netatmo client ID and client secret in this file. We define the port 9210 and assign the IP 192.168.3.5152 of our network as defined in the *networks* (refer above) to it.  

## Netatmo configuration file

In this JSON file named *token.json* you need to define your access and refresh token as well as the expiry date. The expire date will be refreshed automatically so you can use the default date given in the file. If you encounter problems change the date accordingly to yours.  
As *token.json* is included in the *.gitignore* of the root project you will not find it in the *data* directory of this example. Keep in mind that you have to create it accordingly.  

```json
{
    "access_token":"your access token",
    "refresh_token":"your refresh token",
    "expiry":"2024-02-15T18:50:09.209911483Z"
}
```

## Watchtower docker-compose configuration

This is the last configuration step to ensure the consistency of your configuration. It's basically a default configuration. Again the settings are below the *services* section. 

```yml
watchtower:
    # monitor the running container and update them if neccassary
    image: containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ${HOME}/.docker/config.json:/config.json
```

That are all steps to get the system up and running. For the in detail Netatmo configuration you can refer to the documentation of [netatmo-exporter on Github](https://github.com/xperimental/netatmo-exporter). For testing purposes you can also use the dashboard provided there before you start defining your private dashboard. Keep in mind that it takes some time after starting the containers before you see them with valid data in your browser.
*localhost:3000* refers to the Grafana container, *localhost:9210* to the netatmo-exporter container (good to see the available metrics) and *localhost:9090* to the Prometheus container (good to check the connection to the netatmo-exporter).  
