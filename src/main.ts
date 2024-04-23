import * as pageView from "./handlers/page-view";

Bun.serve({
	hostname:
		process.env.NODE_ENV === "production" ? "rando.sahithyan.dev" : "0.0.0.0",
	fetch(req) {
		const url = new URL(req.url);

		switch (url.pathname) {
			case "/page-view":
				return pageView.GET(req);
			default:
				return new Response("Not found", {
					status: 404,
				});
		}
	},
});
