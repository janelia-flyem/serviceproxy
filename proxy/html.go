package proxy

const htmldata=`
<html>
<head>
    <script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.2.5/angular.min.js"></script>
</head>
<body ng-app="myapp">
<div ng-controller="MyController" >
    <button ng-click="myData.doClick(item, $event)">Search Available Services</button>
    <br/>
    <h2>Available Services</h2>
    <div ng-repeat="step in myData.fromServer">
    <div ng-repeat="(key, value) in step">
        <a href="/services/{{key}}/interface">{{key}}</a> (<a href="/services/{{key}}/node">retrieve random node</a>)
    </div>
    </div>
</div>

<script>
    angular.module("myapp", [])
        .controller("MyController", function($scope, $http) {
            $scope.myData = {};
            $scope.myData.doClick = function(item, event) {
                var responsePromise = $http.get("/services");
                responsePromise.success(function(data, status, headers, config) {
                        $scope.myData.fromServer = data.services;
                });
                responsePromise.error(function(data, status, headers, config) {
                    alert("AJAX call failed!");
                });
            }
        } );
</script>
<body ng-app="myapp">
</html>
`


