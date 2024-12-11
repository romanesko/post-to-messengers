# REST service to send XMPP messages


## Prepare and run

### .env file example

```
SERVICE_PORT=8080
XMPP_SERVER=localhost
XMPP_USERNAME=username
XMPP_PASSWORD=password
```

### build

```bash
go build -o out/xmpp-sender
```

### run

```bash
set -a && source .env && ./xmpp-sender
```


### extra

repo contains `Makefile` to build and run the service on amd64 and arm64 platforms

## Usage

### curl

```bash
curl -H 'Content-Type: application/json' \
-d '{"recipient": "test1@localhost","message": "hello there"}' \
'http://localhost:8080/send'
```

### php

```php
<?php

function sendMessage($recipient, $message) {
    $curl = curl_init();

    curl_setopt_array($curl, array(
        CURLOPT_URL => 'http://localhost:8080/send',
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_ENCODING => '',
        CURLOPT_MAXREDIRS => 10,
        CURLOPT_TIMEOUT => 0,
        CURLOPT_FOLLOWLOCATION => true,
        CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
        CURLOPT_CUSTOMREQUEST => 'POST',
        CURLOPT_POSTFIELDS => json_encode(array(
            "recipient" => $recipient,
            "message" => $message
        )),
        CURLOPT_HTTPHEADER => array(
            'Content-Type: application/json'
        ),
    ));

    $response = curl_exec($curl);
    curl_close($curl);

    return $response;
}

/*
// Example usage:

$response = sendMessage("test1@localhost", "hello there");
echo $response;
*/

?>
```