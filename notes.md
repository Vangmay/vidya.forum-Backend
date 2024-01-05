- Username is unique, add uniquiness check in controllers.register (DONE)

- Patch should receive user ID internally instead of as a URL parameter for security reasons, for edit user functionality (DONE)

- Add ability to edit password for user, mabye give them a secret randomly generated code for super authentication (DONE)

- Email is not getting updated (DONE)

- When sending get request to get all posts, every post is not getting preloaded with the comments, that can be a problem. Also change it so that id is in descending order so latest posts are visible first. Fix it in controllers.GetAllPosts. (DONE)

- Add functionality to see if post has been edited or not (DONE)

- Getting most popular posts gives post in ascending order, i need descending (DONE)

- Similarly to the previous issue, getting post by tags give posts in ascending order (DONE)

- Get comments by postId has a bug which causes it to return comments from the wrong post (DONE)

- Deleting comment by comment ID gives unauthorized user idk why
  (Done)
