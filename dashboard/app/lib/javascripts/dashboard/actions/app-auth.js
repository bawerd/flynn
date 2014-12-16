//= require ../dispatcher

(function () {

"use strict";

var Dispatcher = Dashboard.Dispatcher;

Dashboard.Actions.AppAuth = {
	createRelease: function (storeId, release) {
		Dispatcher.handleViewAction({
			name: "APP_ENV:CREATE_RELEASE",
			storeId: storeId,
			release: release
		});
	}
};

})();
