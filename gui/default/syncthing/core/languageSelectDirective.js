angular.module('syncthing.core')
    .directive('languageSelect', ['$document', 'LocaleService', function ($document, LocaleService) {
        'use strict';
        return {
            restrict: 'EA',
            template:
                '<a ng-if="visible" href="" ng-click="toggleLanguageMenu($event)"><span class="fas fa-globe"></span><span class="hidden-xs">&nbsp;{{localesNames[currentLocale] || "English"}}</span> <span class="caret"></span></a>' +
                '<ul ng-if="visible" class="dropdown-menu" ng-show="languageMenuOpen" ng-click="$event.stopPropagation()">' +
                '<li ng-repeat="name in commonLocaleNames" ng-class="{active: localesNamesInv[name]==currentLocale}">' +
                '<a href="" data-ng-click="changeLanguage(localesNamesInv[name], $event)">{{name}}</a>' +
                '</li>' +
                '<li ng-if="otherLocaleNames.length > 0" class="divider"></li>' +
                '<li ng-if="otherLocaleNames.length > 0">' +
                '<a href="#" data-ng-click="toggleMoreLocales($event)">{{showAllLocales ? "收起其他语言" : "更多语言"}}</a>' +
                '</li>' +
                '<li ng-repeat="name in otherLocaleNames" ng-if="showAllLocales" ng-class="{active: localesNamesInv[name]==currentLocale}">' +
                '<a href="" data-ng-click="changeLanguage(localesNamesInv[name], $event)">{{name}}</a>' +
                '</li>' +
                '</ul>',

            link: function ($scope, $element, $attrs) {
                var availableLocales = LocaleService.getAvailableLocales();
                var localeNames = LocaleService.getLocalesDisplayNames();
                var availableLocaleNames = {};
                var englishLanguageAliases = {
                    "zh-CN": "Chinese (Simplified)",
                    "zh-TW": "Chinese (Traditional)",
                    "zh-HK": "Chinese (Hong Kong)",
                };

                // get only locale names that present in available locales
                for (var i = 0; i < availableLocales.length; i++) {
                    var a = availableLocales[i];
                    if (englishLanguageAliases[a]) {
                        availableLocaleNames[a] = englishLanguageAliases[a];
                    } else if (localeNames[a]) {
                        availableLocaleNames[a] = localeNames[a];
                    } else {
                        // show code lang if it is not in the dict
                        availableLocaleNames[a] = '[' + a + ']';
                    }
                }
                $scope.localesNames = availableLocaleNames;

                var invert = function (obj) {
                    var new_obj = {};

                    for (var prop in obj) {
                        if (obj.hasOwnProperty(prop)) {
                            new_obj[obj[prop]] = prop;
                        }
                    }
                    return new_obj;
                };
                $scope.localesNamesInv = invert($scope.localesNames);
                $scope.localesNamesInvKeys = Object.keys($scope.localesNamesInv).sort();
                $scope.showAllLocales = false;
                $scope.languageMenuOpen = false;

                // Keep common languages visible, collapse the rest.
                var commonLocaleCodes = [
                    "en", "zh-CN", "zh-TW", "ja", "ko", "fr", "de", "es", "ru"
                ];
                var commonNameSet = {};
                commonLocaleCodes.forEach(function (code) {
                    if ($scope.localesNames[code]) {
                        commonNameSet[$scope.localesNames[code]] = true;
                    }
                });
                $scope.commonLocaleNames = [];
                $scope.otherLocaleNames = [];
                $scope.localesNamesInvKeys.forEach(function (name) {
                    if (commonNameSet[name]) {
                        $scope.commonLocaleNames.push(name);
                    } else {
                        $scope.otherLocaleNames.push(name);
                    }
                });
                if ($scope.commonLocaleNames.length === 0) {
                    $scope.commonLocaleNames = $scope.localesNamesInvKeys.slice(0, Math.min(6, $scope.localesNamesInvKeys.length));
                    $scope.otherLocaleNames = $scope.localesNamesInvKeys.slice($scope.commonLocaleNames.length);
                }

                $scope.visible = $scope.localesNames && $scope.localesNames['en'];

                // using $watch cause LocaleService.currentLocale will be change after receive async query accepted-languages
                // in LocaleService.readBrowserLocales
                var remove_watch = $scope.$watch(LocaleService.getCurrentLocale, function (newValue) {
                    if (newValue) {
                        $scope.currentLocale = newValue;
                        remove_watch();
                    }
                });

                $scope.changeLanguage = function (locale, $event) {
                    if ($event) {
                        $event.preventDefault();
                        $event.stopPropagation();
                    }
                    LocaleService.useLocale(locale, true);
                    $scope.currentLocale = locale;
                    $scope.languageMenuOpen = false;
                    $scope.showAllLocales = false;
                };

                $scope.toggleMoreLocales = function ($event) {
                    $event.preventDefault();
                    $event.stopPropagation();
                    $scope.showAllLocales = !$scope.showAllLocales;
                };

                $scope.toggleLanguageMenu = function ($event) {
                    $event.preventDefault();
                    $event.stopPropagation();
                    $scope.languageMenuOpen = !$scope.languageMenuOpen;
                };

                // Force submenu open state on host <li> to avoid nested dropdown
                // conflicts with parent menu handlers.
                $scope.$watch('languageMenuOpen', function (open) {
                    $element.toggleClass('open', !!open);
                });

                var closeOnOutsideClick = function () {
                    $scope.$applyAsync(function () {
                        $scope.languageMenuOpen = false;
                        $scope.showAllLocales = false;
                    });
                };
                $document.on('click', closeOnOutsideClick);
                $scope.$on('$destroy', function () {
                    $document.off('click', closeOnOutsideClick);
                });
            }
        };
    }]);
