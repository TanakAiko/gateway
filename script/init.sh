echo "Enter the link of the repository !"
read link
git config --global credential.helper store
git clone $link