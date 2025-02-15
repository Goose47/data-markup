import { Card, Loader } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import { MarkupType, MarkupTypeField, MarkupTypeFull } from "../../utils/types";
import "./MyMarkupType.scss";
import { useEffect, useMemo, useState } from "react";
import { deleteMarkupType, getDetailedMarkupType } from "../../utils/requests";
import { sleep } from "../../utils/utils";
import { MarkupTypeGroup } from "../MarkupTypeGroup/MarkupTypeGroup";
import { PencilToLine, TrashBin } from "@gravity-ui/icons";
import { ButtonWithConfirm } from "../ButtonWithConfirm/ButtonWithConfirm";
import { Link } from "react-router";

const _ = require("lodash");
const b = block("my-markup-type");

export type MarkupTypeIds = "1" | "2" | "3" | "4" | "5";

export const MyMarkupType = ({
  triggerRerender,
  markupType,
}: {
  triggerRerender: () => void;
  markupType: MarkupType;
}) => {
  const [content, setContent] = useState<MarkupTypeFull>();
  const [loading, setLoading] = useState<boolean>();

  useEffect(() => {
    setLoading(true);
    getDetailedMarkupType(markupType.id).then(async (data) => {
      await sleep(500);
      setContent(data);
      setLoading(false);
    });
  }, [markupType.id]);

  const groupedFields = useMemo(() => {
    const result: Record<number, MarkupTypeField[]> = {};
    content?.fields.forEach((field) => {
      if (!result[field.group_id]) {
        result[field.group_id] = [];
      }
      result[field.group_id].push(field);
    });
    return result;
  }, [content]);

  const [values, setValues] = useState<string[][]>([]);

  useEffect(() => {
    if (groupedFields) {
      setValues(
        Object.values(groupedFields).map((value: MarkupTypeField[]) => {
          if (value.length > 0) {
            if (value[0].assessment_type_id === 1) {
              return [String(value[0].id)];
            } else {
              return [];
            }
          }
          return [""];
        })
      );
    }
  }, [groupedFields]);

  const handleUpdateValues = (index: number, newValues: string[]) => {
    const valuesCopy = _.cloneDeep(values);
    valuesCopy[index] = newValues;
    setValues(valuesCopy);
  };

  useEffect(() => {
    console.log(values);
  }, [values]);

  return (
    <div className={b()}>
      <Card className={b("card")}>
        {loading ? (
          <Loader />
        ) : (
          content && (
            <div className={b("wrapper")}>
              <div className={b("edit")}>
                <Link to={`/markup/${markupType.id}`}>
                  <PencilToLine width={20} height={20}></PencilToLine>
                </Link>
              </div>
              <div className={b("delete")}>
                <ButtonWithConfirm
                  handleSubmit={() => {
                    console.log("GAGAG");
                    deleteMarkupType(markupType.id).then(async () => {
                      triggerRerender();
                    });
                  }}
                  confirmText="Вы действительно хотите удалить этот тип разметки?"
                >
                  <TrashBin width={20} height={20}></TrashBin>
                </ButtonWithConfirm>
              </div>
              <h2>{content.name}</h2>
              <div>
                {Object.entries(groupedFields).map(
                  ([_, markupItems], index) => {
                    return (
                      <MarkupTypeGroup
                        value={values[index]}
                        onUpdate={(v) => handleUpdateValues(index, v)}
                        fields={markupItems}
                      />
                    );
                  }
                )}
              </div>
            </div>
          )
        )}
      </Card>
    </div>
  );
};
