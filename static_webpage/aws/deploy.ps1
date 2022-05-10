$BUCKET_NAME="static-webpage-1"
$REGION=$(aws configure get region)
$WEBSITE_URL="http://$BUCKET_NAME.s3-website-$REGION.amazonaws.com"
$HOSTED_ZONE_ID=$(aws route53 list-hosted-zones-by-name --output text --dns-name pkonkol.link | Select-Object -Index 1 | ForEach-Object { $_.split()[2].trimStart("/hostedzone/") })

$Env:AWS_PROFILE="static_webpage"
aws s3 mb s3://$BUCKET_NAME
aws s3 website s3://$BUCKET_NAME/ --index-document index.html --error-document error.html
aws s3 cp ../ s3://$BUCKET_NAME --exclude "aws/*" --recursive
aws s3api put-bucket-policy --bucket static-webpage-1 --policy file://bucket_public_read.json
echo $WEBSITE_URL
aws route53 change-resource-record-sets --hosted-zone-id $HOSTED_ZONE_ID --change-batch file://dns.json