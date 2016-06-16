# aws-coreos-dashboard
AWS CoreOS dashboard

```
env "AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID" \
"AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY" \
"AWS_REGION=eu-west-1" "IP_ADDRESSES=public" \
aws-coreos-dashboard
```

or

```
/usr/bin/docker run -it --rm --name aws-coreos-dashboard -p 8080:8080 \
-e "AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID" \
-e "AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY" \
-e "AWS_REGION=eu-west-1" \
-e "IP_ADDRESSES=private" \
aws-coreos-dashboard
```
