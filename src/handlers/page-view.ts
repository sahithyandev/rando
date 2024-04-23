// Imports the Google Analytics Data API client library.
import { BetaAnalyticsDataClient, type protos } from "@google-analytics/data";
import { join } from "path";
import { readFileStructure } from "../utils";

/**
 * TODO(developer): Uncomment this variable and replace with your
 *   Google Analytics 4 property ID before running the sample.
 */
const propertyId = "356543625";

console.log("process.cwd", process.cwd());
console.log(
	"readdir",
	readFileStructure(process.cwd(), ["node_modules", ".git"])
);
process.env.GOOGLE_APPLICATION_CREDENTIALS = join(
	process.cwd(),
	"./blog-pageviews-1713811580440-51ea0fe058d1.json"
);
console.log(
	"GOOGLE_APPLICATION_CREDENTIALS",
	process.env.GOOGLE_APPLICATION_CREDENTIALS
);

// Using a default constructor instructs the client to use the credentials
// specified in GOOGLE_APPLICATION_CREDENTIALS environment variable.
const analyticsDataClient = new BetaAnalyticsDataClient();

export async function GET(request: Request) {
	const url = new URL(request.url);
	const pagePath = url.searchParams.get("path");
	console.log("pagePath", pagePath);

	const runReportOptions: protos.google.analytics.data.v1beta.IRunReportRequest =
		{
			property: `properties/${propertyId}`,
			dimensions: [{ name: "pagePath" }],
			metrics: [{ name: "screenPageViews" }],
			dateRanges: [{ startDate: "2022-03-05", endDate: "yesterday" }],
		};
	if (pagePath) {
		runReportOptions.dimensionFilter = {
			filter: {
				fieldName: "pagePath",
				stringFilter: { matchType: "EXACT", value: pagePath },
			},
		};
	}

	const responseArr = await analyticsDataClient
		.runReport(runReportOptions)
		.catch(console.error);
	const response = responseArr === undefined ? undefined : responseArr[0];
	console.log("response", response);
	const rows = response?.rows;
	if (!rows) {
		console.log(response);
		return new Response("Not available now", {
			status: 500,
		});
	}

	const map: Record<string, string> = {};
	for (const row of rows) {
		if (!row.dimensionValues || !row.metricValues) {
			continue;
		}
		const key = row.dimensionValues[0].value;
		const value = row.metricValues[0].value;
		console.log(key, value, !pagePath);
		if (!pagePath && typeof key === "string" && typeof value === "string") {
			map[key] = value;
		} else if (pagePath === key) {
			return new Response(value || "0", {
				status: 200,
			});
		}
	}

	return new Response(
		JSON.stringify({
			pages: map,
			dataRetrievedOn: new Date().toISOString(),
		}),
		{
			status: 200,
			headers: {
				"Content-Type": "application/json",
				"Cache-Control": "public, max-age=",
			},
		}
	);
}
