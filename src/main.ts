import * as pageView from "./handlers/page-view";

Bun.serve({
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
