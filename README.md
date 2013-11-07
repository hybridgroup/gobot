## Cylon Website

This site is build using [Middleman](http://middlemanapp.com/getting-started/)  
  
To run locally:  

      bundle install
      bundle exec middleman

### Deploy

[middleman-gh-pages](https://github.com/neo/middleman-gh-pages) gem is being used to build the webpage and deploy to gh-pages branch.  

For deploying the webpage, your must be in 'Cylon.js.io' branch and run the following command:

      rake publish

