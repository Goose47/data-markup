import { Alert, Button, Flex, Spin } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import {
  assessmentNext,
  assessmentStore,
  assessmentUpdate,
  batchFind,
  getAvailableMarkupTypes,
  getBatchMarkupData,
  getDetailedMarkupType,
} from "../../utils/requests";
import { AssessmentNext, MarkupType, MarkupTypeFull } from "../../utils/types";
import "./Assessment.scss";
import { useEffect, useState } from "react";
import {
  FieldValue,
  MyMarkupType,
} from "../../components/MyMarkupType/MyMarkupType";
import { MarkupData } from "../../components/MarkupData/MarkupData";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";

const b = block("assessment");

const valueToFields = (values: FieldValue[][]) => {
  return {
    fields: values.flatMap((field) => {
      return field.map((value) =>
        value.assessment_type_id === 5
          ? {
              text: value.value,
              markup_type_field_id: value.fieldIdFuckBackend ?? 0,
            }
          : {
              text: null,
              markup_type_field_id: parseInt(value.value),
            }
      );
    }),
  };
};

export type BatchType = "single" | "compare";

export type AssessmentData = AssessmentNext & { batchType: BatchType };

export const Assessment = ({
  markupId,
  isAdmin,
  triggerRerender,
}: {
  markupId?: number;
  isAdmin?: boolean;
  triggerRerender?: () => void;
}) => {
  const [currentAssessment, setCurrentAssessment] = useState<
    AssessmentData | "loading" | "error" | null
  >(null);
  const [currentValue, setCurrentValue] = useState<FieldValue[][]>([[]]);

  const hasCurrentAssessment =
    currentAssessment !== null &&
    currentAssessment !== "loading" &&
    currentAssessment !== "error";

  useEffect(() => {
    const fetchAssessmentNext = async () => {
      if (currentAssessment !== null) {
        return;
      }

      const assessment: AssessmentNext = await assessmentNext();

      if (!assessment?.markup_type?.batch_id) {
        setCurrentAssessment("error");
        return;
      }

      const batchType: { type_id: number } = await batchFind(
        assessment.markup_type.batch_id
      );

      setCurrentAssessment({
        ...assessment,
        batchType: batchType.type_id === 1 ? "single" : "compare",
      });
    };

    const fetchAssessmentByMarkup = async (markupId: number) => {
      if (currentAssessment !== null) {
        return;
      }

      const rawAssessment = await getBatchMarkupData(markupId);

      const batchType: { type_id: number } = await batchFind(
        rawAssessment.batch_id
      );

      getAvailableMarkupTypes(parseInt(rawAssessment.batch_id)).then(
        (data: MarkupType[]) => {
          getDetailedMarkupType(data[data.length - 1].id).then(
            (data: MarkupTypeFull) => {
              const assessment: AssessmentNext = {
                ...rawAssessment,
                markup_type: data,
              };
              setCurrentAssessment({
                ...assessment,
                batchType: batchType.type_id === 1 ? "single" : "compare",
              });
            }
          );
        }
      );
    };

    if (markupId) {
      fetchAssessmentByMarkup(markupId);
    } else {
      fetchAssessmentNext();
    }
  }, [currentAssessment, markupId]);

  return (
    <div className={b("wrapper")}>
      {!hasCurrentAssessment ? (
        <>
          {currentAssessment === "loading" && <Spin />}
          {currentAssessment === "error" && (
            <Alert title="Кончились задания 🤯" theme="danger" />
          )}
        </>
      ) : (
        <Flex direction="column" gap={4}>
          <MarkupData assessment={currentAssessment} />
          <MyMarkupType
            isAdmin={false}
            triggerRerender={() => {}}
            markupType={currentAssessment.markup_type}
            onUpdateValue={setCurrentValue}
          />

          <Button
            view="action"
            onClick={async () => {
              setCurrentAssessment("loading");
              if (isAdmin && markupId) {
                await assessmentStore({
                  markup_id: markupId,
                  ...valueToFields(currentValue),
                });
              } else {
                await assessmentUpdate(
                  currentAssessment.assessment_id ?? -1,
                  valueToFields(currentValue)
                );
              }
              if (triggerRerender) {
                triggerRerender();
              }
              toaster.add({
                name: "Отправлено",
                content: "+0.03 ₽",
                theme: "success",
              });
              setCurrentAssessment(null);
            }}
          >
            {isAdmin ? "Установить эталонную оценку" : "Отправить"}
          </Button>
        </Flex>
      )}
    </div>
  );
};
