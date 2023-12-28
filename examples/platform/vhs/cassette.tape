# To run:
# cd example/platform
# vhs ./vhs/cassette.tape

Output ./vhs/platform_usage.gif
Require terraform
Require terrarium

Sleep 2s
Type "pwd"
Enter
Type "ls -la"
Enter
Sleep 2s
Type "ls -la mock-modules"
Enter
Sleep 2s
Type "terraform init && terraform plan"
Enter
Sleep 4s
Type "rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup"
Enter
Sleep 2s
Type "terrarium platform lint"
Enter
Sleep 4s
Type "terrarium generate -c dev -a ../apps/voting-be -a ../apps/voting-fe -a ../apps/voting-worker"
Enter
Sleep 4s
Type "ls -la .terrarium"
Enter
Sleep 2s
Type "cd .terrarium"
Enter
Sleep 2s
Type "terraform init && terraform plan && terraform apply"
Enter
Sleep 4s
Type "yes"
Enter
Sleep 2s
Type "ls -la"
Enter
Sleep 2s
Type "terraform output -json | mustache app_voting_worker.env.mustache"
Enter
Sleep 2s
Type "terraform output -json | mustache app_voting_be.env.mustache"
Enter
Sleep 2s
Type "terraform output -json | mustache app_voting_fe.env.mustache"
Enter
Sleep 2s
Type "cd .. && rm -rf .terrarium"
Enter
Sleep 2s
