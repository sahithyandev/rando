// Imports the Google Analytics Data API client library.
import { BetaAnalyticsDataClient, type protos } from "@google-analytics/data";
import { join } from "path";
import merge from "deepmerge";

const RESPONSE_MAX_AGE = 108000; // 30 hours in seconds

/**
 * TODO(developer): Uncomment this variable and replace with your
 *   Google Analytics 4 property ID before running the sample.
 */
const propertyId = "356543625";

const propertyIdMap = Object.freeze({
	"sahithyan.dev": "356543625",
	"kalvi.lk": "398313900",
});

function isPropertyAvailableForDomain(
	domain: string
): domain is keyof typeof propertyIdMap {
	return domain in propertyIdMap;
}

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

const headers = {
	"Content-Type": "application/json",
	"Cache-Control": `public, max-age=${RESPONSE_MAX_AGE}`,
};

export async function GET(
	request: Request,
	responseInitOptions: ResponseInit = {}
) {
	const url = new URL(request.url);
	const pagePath = url.searchParams.get("path");

	let targetDomain = url.searchParams.get("domain");
	if (!targetDomain) {
		request.headers.forEach((value, key) => {
			console.log("header", key, value);
		});
		const originHeader = request.headers.get("Origin");
		if (originHeader) {
			targetDomain = originHeader.split(":")[0].split("://")[1];
		}
	}

	console.log("getting page views for", targetDomain, "on", pagePath);
	if (!targetDomain || !isPropertyAvailableForDomain(targetDomain)) {
		return new Response(
			"Unsupported domain",
			merge<ResponseInit>(
				{
					status: 400,
				},
				responseInitOptions
			)
		);
	}

	const runReportOptions: protos.google.analytics.data.v1beta.IRunReportRequest =
		{
			property: `properties/${propertyIdMap[targetDomain]}`,
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
	const rows = response?.rows;
	if (!rows) {
		console.log(response);
		return new Response(
			"Not available now",
			merge<ResponseInit>(
				{
					status: 500,
				},
				responseInitOptions
			)
		);
	}

	const map: Record<string, string> = {};
	const _responseOptions = merge<ResponseInit>(
		{
			status: 200,
			headers,
		},
		responseInitOptions
	);
	console.log("merged response options", _responseOptions);

	for (const row of rows) {
		if (!row.dimensionValues || !row.metricValues) {
			continue;
		}
		const key = row.dimensionValues[0].value;
		const value = row.metricValues[0].value;
		if (!pagePath && typeof key === "string" && typeof value === "string") {
			map[key] = value;
		} else if (pagePath === key) {
			return new Response(
				JSON.stringify({
					views: value || "0",
					dateRetrievedOn: new Date().toISOString(),
				}),
				_responseOptions
			);
		}
	}

	return new Response(
		JSON.stringify({
			pages: map,
			dataRetrievedOn: new Date().toISOString(),
		}),
		_responseOptions
	);
}
