This doc outlines the spec for an application I'd like you to create.

Overview:

The application is a blog, written in Golang and React, and backed by a mysql database. User accounts are created via OIDC. Users can comment on posts, react to comments, and share posts. 

Security foreward:

Do *not* embed any DB credentials or API keys/secrets anywhere in the source of the project. In spaces where a key is needed, create an entry for it and supply a dummy value into `.env.local` at the project root, and then reference that environment variable instead.

Infrastructure:

The application will consist of three separate docker containers -- one for the back-end, one for the front-end, and lastly a DB container. There will be a project root-level docker compose file orchestrating the containers.

Features:

- Main page: The main page hosts a list of blog posts, along with their post images (if applicable) and created dates. Each listing links to its associated post.

- Blog posts: Blog posts will be uploaded directly by the admin user. No admin client necessary. columns associated with posts will include: Content, Date created, Date updated, arbitrary string tags, optional post image, slug, title, excerpt, canonical URL, and meta description. When displaying blog posts, display the post image as a hero image, followed by the title, created date (updated date in parenthetical if applicable), and then the post content. The site should also include a minimal progress bar flush with the bottom of the page that indicates how far along the blog post the user is

- Comments: Comments can be submitted on a specific post and all comments on posts are displayed vertically in order of submission at the bottom of the post. Each comment has an associated user and can support users reacting with any emoji in the standard unicode set. Comments are displayed with user display name at the top, comment content, then reactions and reaction controls. Comments do not support direct replies. When a user submits a comment, the data should be sanitized into html-safe characters before storing into the comments table. Comment reactions should be stored in their own table, including the unicode character of the emoji, the associated comment Id, a count of the number of reactions of that emoji type. A new reaction should query the table for that specific reaction emoji/comment combo, adding 1 to the row's count if it exists, and inserting a new row w/ count 1 if not.

- User Accounts: User accounts are managed via OIDC. The blog supports five specific OIDC providers -- Google, Apple, Facebook, LinkedIn, and Github. The blog requests minimal claims necessary to get the user a display name for their comments. User accounts are represented in the database by a Users table, which contains the display name and Id, and a User identities table that contains the association with the Users table as well as the user's linked information from OIDC. In asking for the claims, we specifically highlight that the purpose for the information is purely to associate users with comments, and that we will never share information provided with us with a third party. Client IDs and Secrets will be provided later, so create placeholders to be filled by environment variables where necessary.

- Sharing: At the bottom of the post, before the comments section, are links to share to facebook, twitter, linkedin, email, and direct link copy. Sharing should include OpenGraph meta information associated with the post.