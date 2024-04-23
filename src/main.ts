import * as pageView from "./handlers/page-view";

const CORS_HEADERS = {
	headers: {
		"Access-Control-Allow-Origin": "*",
		"Access-Control-Allow-Methods": "GET,HEAD,PUT,PATCH,POST,DELETE",
		"Access-Control-Allow-Headers": "Content-Type",
	},
};

Bun.serve({
	fetch(request) {
		const url = new URL(request.url);

		// Handle CORS preflight requests
		if (request.method === "OPTIONS") {
			return new Response("Departed", CORS_HEADERS);
		}

		switch (url.pathname) {
			case "/page-view":
				return pageView.GET(request);
			default:
				return new Response("Not found", {
					status: 404,
				});
		}
	},
});
