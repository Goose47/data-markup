import { Flex, Link, Text } from "@gravity-ui/uikit";
import { AssessmentData, BatchType } from "../../pages/Assessment/Assessment";
import { block } from "../../utils/block";
import "./MarkupData.scss";

const b = block("markup-data");

type MarkupFieldEntity = { key: string; value: string; type: "text" | "img" | "url" };

/**
 * keys are parsed by _postfix: _text, _img, _url
 * group is parsed by some_nameN_postfix where 0 <= N <= 9
 */
const toParsedFields = (data: any, type: BatchType) => {
    const parseKey = (key: string): {
        key: string;
        type: "text" | "img" | "url";
        group: number | null;
    } => {
        const splitted = key.split("_");

        if (splitted.length === 1) {
            const group = parseInt(splitted.at(-1)?.at(-1) ?? "");

            return {
                key,
                group: !isNaN(group) ? group : null,
                type: "text",
            };
        }

        const postfix = splitted.at(-1);
        const group = parseInt(splitted.at(-2)?.at(-1) ?? "");

        let parsed;
        switch (postfix) {
            case "text":
                parsed = {
                    key: splitted.slice(0, -1).join("_"),
                    type: "text" as const,
                };
                break;
            case "img":
                parsed = {
                    key: splitted.slice(0, -1).join("_"),
                    type: "img" as const,
                };
                break;
            case "url":
                parsed = {
                    key: splitted.slice(0, -1).join("_"),
                    type: "url" as const,
                };
                break;
            default:
                parsed = {
                    group: !isNaN(parseInt(key.at(-1) ?? ""))
                        ? parseInt(key.at(-1) ?? "")
                        : null,
                    key,
                    type: "text" as const,
                };
                break;
        }

        return { group: type === "compare" && !isNaN(group) ? group : null, ...parsed };
    };

    const fields = Object.entries(data).map(([key, value]) => ({
        ...parseKey(key),
        value: value as string,
    }));

    return fields;
};

const MarkupField = ({ field }: { field: MarkupFieldEntity }) => {
    let value;

    switch (field.type) {
        case "img":
            value = <img src={field.value} alt={field.key} className={b("image")} />;
            break;
        case "url":
            value = <Link href={field.value}>{field.value}</Link>;
            break;
        case "text":
        default:
            value = <Text variant="body-1">{field.value}</Text>;
            break;
    }

    return value;
};

export const MarkupData = ({ assessment }: { assessment: AssessmentData }) => {
    const data = JSON.parse(assessment.data);
    const fields = toParsedFields(data, assessment.batchType);

    if (assessment.batchType === "compare") {
        // Group fields by group number; ungrouped fields go to a special key.
        const UNGROUPED = "Ungrouped";
        const groupedFields = Object.groupBy(fields, (field) => field.group ?? UNGROUPED);
        const ungrouped = groupedFields[UNGROUPED] ?? [];
        delete groupedFields[UNGROUPED];

        const groupKeys = Object.keys(groupedFields);
        if (groupKeys.length === 0) {
            // Fallback if no groups exist.
            return (
                <div className={b("wrapper")}>
                    <Flex direction="column" gap={4}>
                        {fields.map((field, i) => (
                            <MarkupField key={i} field={field} />
                        ))}
                    </Flex>
                </div>
            );
        }

        // Use the first available group to establish ordering for group-specific rows.
        const firstGroup = Number(groupKeys[0]);
        const groupSpecificCount = groupedFields[firstGroup]?.length ?? 0;
        // Total rows: header row + (ungrouped fields + group-specific fields)
        const totalRows = 1 + ungrouped.length + groupSpecificCount;

        // Build left-column cells:
        // First cell is the header label (e.g. "Группа1"), then ungrouped field keys, then group-specific field keys.
        const leftColumnCells: React.ReactNode[] = [];
        leftColumnCells.push("Группа1");
        ungrouped.forEach((field) => leftColumnCells.push(field.key));
        for (let i = 0; i < groupSpecificCount; i++) {
            leftColumnCells.push(groupedFields[firstGroup]?.[i].key);
        }

        // Build each group's column.
        const groupColumns: React.ReactNode[][] = [];
        groupKeys.forEach((group) => {
            const groupArr = groupedFields[group as unknown as number] ?? [];
            const colCells: React.ReactNode[] = [];
            colCells.push(<Text variant="header-1" key={`header-${group}`}>{group}</Text>);
            ungrouped.forEach((field) =>
                colCells.push(
                    <MarkupField key={`ungrouped-${group}-${field.key}`} field={field} />
                )
            );
            groupArr.forEach((field, index) =>
                colCells.push(
                    <MarkupField key={`group-${group}-${field.key}-${index}`} field={field} />
                )
            );
            groupColumns.push(colCells);
        });

        // Total columns: left column + one column per group.
        const totalCols = 1 + groupColumns.length;
        const gridCells: React.ReactNode[] = [];

        // Flatten the grid cells (row by row).
        for (let row = 0; row < totalRows; row++) {
            // Left column cell.
            gridCells.push(
                <div key={`cell-${row}-0`} className={b("grid-cell", "label")}>
                    <Text variant="header-1">{leftColumnCells[row]}</Text>
                </div>
            );
            // Each group's cell.
            for (let col = 0; col < groupColumns.length; col++) {
                gridCells.push(
                    <div key={`cell-${row}-${col + 1}`} className={b("grid-cell")}>
                        {groupColumns[col][row]}
                    </div>
                );
            }
        }

        return (
            <div className={b("wrapper")}>
                <div
                    className={b("grid-container")}
                    style={{
                        display: "grid",
                        gridTemplateColumns: `repeat(${totalCols}, 1fr)`,
                        gap: "8px",
                    }}
                >
                    {gridCells}
                </div>
            </div>
        );
    }

    // For non-"compare" batch types, render as before.
    return (
        <div className={b("wrapper")}>
            <Flex direction="column" gap={4}>
                {fields.map((field, index) => (
                    <MarkupField key={index} field={field} />
                ))}
            </Flex>
        </div>
    );
};
