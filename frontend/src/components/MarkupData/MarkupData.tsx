import { Flex, Link, Text } from "@gravity-ui/uikit";
import { AssessmentData, BatchType } from "../../pages/Assessment/Assessment";
import { block } from "../../utils/block";
import "./MarkupData.scss";

const b = block("markup-data");

type MarkupFieldEntity = {
  key: string;
  value: string;
  type: "text" | "img" | "url";
};

/**
 * keys are parsed by _postfix: _text, _img, _url
 * group is parsed by some_nameN_postfix where 0 <= N <= 9
 */
const toParsedFields = (data: any, type: BatchType) => {
  const parseKey = (
    key: string
  ): {
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

    return {
      group: type === "compare" && !isNaN(group) ? group : null,
      ...parsed,
    };
  };

  const fields = Object.entries(data).map(([key, value]) => ({
    ...parseKey(key.trim()),
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
      value = (
        <Link target="_blank" href={field.value}>
          {field.value}
        </Link>
      );
      break;
    case "text":
      value = <Text variant="body-1">{field.value}</Text>;
      break;
    default:
      value = <Text variant="body-1">{field.value}</Text>;
      break;
  }

  return (
    <Flex gap={4} justifyContent={"space-between"}>
      {/* <Text variant="header-1">{field.key}</Text> */}
      {value}
    </Flex>
  );
};

export const MarkupData = ({ assessment }: { assessment: AssessmentData }) => {
  const data = JSON.parse(assessment.data);

  const fields = toParsedFields(data, assessment.batchType);

  const UNGROOUPED = "Ungrouped";

  if (assessment.batchType === "compare") {
    const groupedFields = Object.groupBy(
      fields,
      (field) => field.group ?? UNGROOUPED
    );
    console.log(groupedFields);
    return (
      <div className={b("wrapper")}>
        <div className={b("compare")}>
          <table>
            <thead>
              <th></th>
              {Object.keys(groupedFields).map((el) => (
                <th>{el}</th>
              ))}
            </thead>
            <tbody>
              {[
                ...new Array(
                  Math.max(
                    ...Object.values(groupedFields).map((el) => el?.length ?? 0)
                  )
                ),
              ].map((_, index) => (
                <tr>
                  <td>
                    {(Object.values(groupedFields)[0] as any)[
                      index
                    ]?.key?.replace(/\d+$/g, "")}
                  </td>
                  {Object.keys(groupedFields).map((key) => {
                    const x = (
                      groupedFields[key as number | "Ungrouped"] as any
                    )[index];
                    if (x) {
                      return (
                        <td className={b(x.type)}>
                          {<MarkupField field={x} />}
                        </td>
                      );
                    }
                    return <td></td>;
                  })}
                </tr>
              ))}
            </tbody>
          </table>
          {/* {Object.entries(groupedFields).map(([group, fields]) => {
            return (
              <div>
                <Text variant="header-1">{group}</Text>
                {fields?.map((field) => <MarkupField field={field} />)}
              </div>
            );
          })} */}
        </div>
      </div>
    );
  }

  return (
    <div className={b("wrapper")}>
      <Flex direction={"column"} gap={4}>
        {fields.map((field) => (
          <MarkupField field={field} />
        ))}
      </Flex>
    </div>
  );
};
