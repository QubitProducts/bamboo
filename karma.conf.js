module.exports = function(karma, specificOptions) {
  karma.set({

    // base path, that will be used to resolve files and exclude
    basePath: "",

    // frameworks to use
    frameworks: ["browserify", "mocha"],
    client: {
      mocha: {
        ui: 'bdd'
      }
    },

    // list of files / patterns to load in the browser
    files: [
      "webapp/dist/main-libs.js",
      "node_modules/expect.js/index.js",
      {pattern: "webapp/**/*_test.js", included: false}
      
    ],
    
    // enable / disable watching file and executing tests whenever any file changes
    autoWatch: true,

    browserify: {
      watch: true
    },

    preprocessors: {
      "webapp/**/*_test.js": ["browserify"]
    },

    // list of files to exclude
    exclude: [
      // "node_modules/**/*"
    ],

    // test results reporter to use
    // possible values: 'dots', 'progress', 'junit', 'growl', 'coverage'
    reporters: ['progress'],

    // web server port
    port: 9876,

    // cli runner port
    runnerPort: 9100,

    // enable / disable colors in the output (reporters and logs)
    colors: true,

    // level of logging
    // possible values: config.LOG_DISABLE || config.LOG_ERROR || config.LOG_WARN || config.LOG_INFO || config.LOG_DEBUG
    logLevel: karma.LOG_DEBUG,

    // Start these browsers, currently available:
    // - Chrome
    // - ChromeCanary
    // - Firefox
    // - Opera
    // - Safari (only Mac)
    // - PhantomJS
    // - IE (only Windows)
    browsers: ['Chrome'],

    // If browser does not capture in given timeout [ms], kill it
    captureTimeout: 60000,

    // Continuous Integration mode
    // if true, it capture browsers, run tests and exit
    singleRun: false

  });
};