### Improvements required

- Sending a delete request while not logged in leads to
  "panic: runtime error: invalid memory address or nil pointer dereference"
  (Perhaps its okay because delete button won't be accessible while the user is not logged in)

- The email and username are all case-sensitive, this can be a problem (FIXED)

- Another user is able to delete a post (FIXED)

### Good stuff

- Authentication works well 👍

Todo

- Complete CHECK 3
- Make code prettier
- Document API
