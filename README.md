## Gobot Website

This site is build using [Middleman](http://middlemanapp.com/basics/getting-started/)  
  
To run locally:  

      bundle install
      bundle exec middleman

### Deploy

[middleman-gh-pages](https://github.com/neo/middleman-gh-pages) gem is being used to build the webpage and deploy to gh-pages branch.  

For deploying the webpage, your must be in 'gobot.io' branch and run the following command:

      rake publish

You must not have any uncomitted or untracked files in the site dirs, or the publish operation will fail with a message such as `Directory not clean`.

If the publish fails, you might need to remove the `build` dir before trying to run `rake publish` again.
