import { Flex, Link, Text } from "@gravity-ui/uikit";
import { AssessmentData, BatchType } from "../../pages/Assessment/Assessment";
import { block } from "../../utils/block";
import "./MarkupData.scss";

const b = block("markup-data");

type MarkupFieldEntity = { key: string, value: string, type: "text" | "img" | "url" };

/**
 * keys are parsed by _postfix: _text, _img, _url
 * group is parsed by some_nameN_postfix where 0 <= N <= 9
 */
const toParsedFields = (data: any, type: BatchType ) => {
    const parseKey = (key: string): {
        key: string;
        type: "text" | "img" | "url";
        group: number | null;
    } => {
        const splitted = key.split("_");

        if (splitted.length === 1) {
            return {
                key,
                group: null,
                type: "text",
            }
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
                    key,
                    type: "text" as const,
                };
                break;
        }

        return { ...parsed, group: type === "compare" && !isNaN(group) ? group : null, };
    }

    const fields = Object.entries(data).map(([key, value]) => ({
        ...parseKey(key),
        value: value as string,
    }))

    return fields;
}

const MarkupField = ({ field }: { field: MarkupFieldEntity}) => {
    let value;

    switch (field.type) {
        case "img":
            value = <img src={field.value} alt={field.key} className={b("image")} />;
            break;
        case "url":
            value = <Link href={field.value}>{field.value}</Link>;
            break;
        case "text":
            value = <Text variant="body-1">{field.value}</Text>;
            break;
        default:
            value = <Text variant="body-1">{field.value}</Text>;
            break;
    }

    return <Flex gap={4} justifyContent={"space-between"}>
        <Text variant="header-1">{field.key}</Text>
        {value}
    </Flex>
}

export const MarkupData = ({ assessment }: { assessment: AssessmentData }) => {
    const data = JSON.parse(assessment.data);

    const fields = toParsedFields(data, assessment.batchType);

    if (assessment.batchType === "compare") {
        fields.reduce((acc, field) => {
            const group = acc[field.group ?? 10];

            if (!group) {
                acc[field.group ?? 10] = [field];
            } else {
                group.push(field);
            }

            return acc;
        }, [] as unknown as [MarkupFieldEntity[]]);
    }

    return (<div className={b("wrapper")}>
        <Flex direction={"column"}>
            {fields.map((field) => <MarkupField field={field}/>)}
        </Flex>
    </div>)
};
