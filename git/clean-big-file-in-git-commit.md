# clean big file in git commit

定位

```sh
git rev-list --objects --all | grep "$(git verify-pack -v .git/objects/pack/*.idx | sort -k 3 -n | tail -5 | awk '{print$1}')"
```

重写commit

```sh
git filter-branch -f --prune-empty --index-filter 'git rm -rf --cached --ignore-unmatch ${your-file-name}' --tag-name-filter cat -- --all
```

提交

```sh
git push origin --force --all
```