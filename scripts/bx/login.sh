#bx login --sso
bx login -c 8cb197d2d1424b5195f365ec6ca188cc -r us-south --sso
bx cr login
bx cs cluster config --cluster KN1
bx target -g Default
kubectl config use-context guestbook-context
