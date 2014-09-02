"use strict";

var _ = require("lodash");

// Common Directories
var distDir = "webapp/dist";
var sassDir = "webapp/sass";
var appDir  = "webapp/app";

module.exports = function (grunt) {
  var pkg = grunt.file.readJSON('package.json');
  var commonShims  = pkg.shims;
  var libSource    = _.pluck(commonShims, 'path');
  var libAlias     = _.map(commonShims, function (shim, key) { return shim.path + ":" + key; });

  console.log("Common Shims:");
  _.map(libAlias, function (lib) { console.log(lib); });

  grunt.initConfig({
    pkg: pkg,

    distDir: distDir,
    sassDir: sassDir,
    appDir: appDir,

    browserify: {

      libs: {
        src: libSource,
        dest: "<%= distDir %>/main-libs.js",
        options: {
          alias: libAlias,
          shim: commonShims,
          transform: ["uglifyify"]
        }
      },

      app: {
        src: ["<%= appDir %>/**/*.js"],
        dest: "<%= distDir %>/main-app.js",
        options: {
          alias: libAlias,
          external: libSource,
          transform: ["partialify", "uglifyify"],
          bundleOptions: {
            debug: false
          }
        }
      },

      all: {
        options: {
          shim: commonShims,
          alias: libAlias,
          transform: ["partialify"]
        },
        src: ["<%= appDir %>/**/*.js"],
        dest: "<%= distDir %>/main-bundle.js"
      }
    },
    sass: {
      dist: {
        options: {                       // Target options
          style: 'compressed'
        },
        files: {
          '<%= distDir %>/main.css' : '<%= sassDir %>/main.scss'
        }
      }
    },
    watch: {
      script: {
        files: ['<%= appDir %>/**/*.js', '<%= appDir %>/**/*.html'],
        tasks: ['browserify:app'],
        options: {
          spawn: false
        }
      },
      sass: {
        files: '<%= sassDir %>/**/*.scss',
        tasks: ['sass']
      }
    }
  });

  grunt.loadNpmTasks('grunt-browserify');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-sass');
  grunt.registerTask('default', ['browserify']);
}