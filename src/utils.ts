import { readdirSync } from "fs";
import { join } from "path";

export function readFileStructure(
	dir: string,
	ignoreDir: string[] | undefined = undefined
) {
	const files = readdirSync(dir, { withFileTypes: true });
	ignoreDir = ignoreDir || [];
	const structured: string[] = [];

	files.forEach((file) => {
		const name = file.name;
		if (file.isFile()) {
			structured.push(join(dir, name));
		} else if (file.isDirectory() && !ignoreDir.includes(name)) {
			readFileStructure(join(dir, name), ignoreDir).forEach((_file) => {
				structured.push(_file);
			});
		}
	});

	return structured;
}
