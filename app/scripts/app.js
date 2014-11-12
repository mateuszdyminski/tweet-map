'use strict';

var app = angular.module('tweetSetApp', ['ngRoute', 'ui.bootstrap', 'leaflet-directive'])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
  });

app.factory('Twitter', function($http, $timeout) {
    var ws = new WebSocket('ws://localhost:9000/wsapi/ws');
    var twitterService = {
      tweets: [],
      requests: [],
      query: function (query, callback) {
        $http({method: 'POST', url: '/restapi/tweets', params: {query: query}}).
          success(function () {
            callback(query);
          }).error(function (err) {
            console.log(err);
          });
      }
    };

    ws.onmessage = function(event) {
      $timeout(function() {
        twitterService.tweets.push(JSON.parse(event.data));
        twitterService.tweets = twitterService.tweets;
      });
    };

    return twitterService;
  });

app.controller('Search', function($scope, $http, $timeout, Twitter) {
  $scope.alerts = [];
  $scope.search = function() {
    Twitter.query($scope.query, $scope.onQueryRequsted);
  };

  $scope.onQueryRequsted = function(query) {
    $scope.alerts = [];
    $scope.alerts.push({msg: 'Request for tweets containing string: ' + query + ' send. Waiting for results.'});
    $scope.query = '';
  };

  $scope.closeAlert = function(index) {
    $scope.alerts.splice(index, 1);
  };
});

app.controller('Tweets', function($scope, $http, $timeout, Twitter) {
  $scope.tweets = [];
  $scope.markers = [];

  $scope.$watch(
    function() {
      return Twitter.tweets;
    },
    function(tweets) {
      $scope.tweets = tweets;

      $scope.markers = tweets.map(function(tweet) {
        return {
          lng: tweet.coordinates.lng,
          lat: tweet.coordinates.lat,
          message: tweet.text,
          focus: true
        };
      });
    }, true);
});
