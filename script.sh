wget https://www.emqx.com/en/downloads/broker/5.7.0/emqx-5.7.0-ubuntu22.04-amd64.tar.gz
mkdir -p emqx && tar -zxvf emqx-5.7.0-ubuntu24.04-amd64.tar.gz -C emqx
ehco '
api_key = {
  bootstrap_file = "etc/default_api_key.conf"
}'> ~/emqx/etc/emqx.conf
echo 'my-app:AAA4A275-BEEC-4AF8-B70B-DAAC0341F8EB'> ~/emqx/etc/default_api_key.conf
./emqx/bin/emqx start
