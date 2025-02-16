import { Alert, Button, Flex, Spin } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import { assessmentNext, assessmentUpdate, batchFind } from "../../utils/requests";
import { AssessmentNext } from "../../utils/types";
import "./Assessment.scss";
import { useEffect, useState } from "react";
import { FieldValue, MyMarkupType } from "../../components/MyMarkupType/MyMarkupType";
import { MarkupData } from "../../components/MarkupData/MarkupData";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";

const b = block("assessment");

const valueToFields = (values: FieldValue[][]) => {
    return {
        fields: values.flatMap(field => {
            return field.map(value => value.assessment_type_id === 5 ? ({
                text: value.value,
                markup_type_field_id: value.fieldIdFuckBackend ?? 0,
            }) : ({
                text: null,
                markup_type_field_id: parseInt(value.value),
            }));
        })
    }
};

export type BatchType = "single" | "compare";

export type AssessmentData = AssessmentNext & { batchType: BatchType };

export const Assessment = () => {
    const [currentAssessment, setCurrentAssessment] = useState<AssessmentData | "loading" | "error" | null>(null);
    const [currentValue, setCurrentValue] = useState<FieldValue[][]>([[]]);

    const hasCurrentAssessment = currentAssessment !== null && currentAssessment !== "loading" && currentAssessment !== "error";

    useEffect(() => {
        const fetchAssessment = async () => {
            if (currentAssessment !== null) {
                return;
            }

            const assessment: AssessmentNext = await assessmentNext();

            if (!assessment?.markup_type?.batch_id) {
                setCurrentAssessment("error");
                return;
            }

            const batchType: { type_id: number } = await batchFind(assessment.markup_type.batch_id);

            setCurrentAssessment({ ...assessment, batchType: batchType.type_id === 1 ? "single" : "compare"});
        }

        fetchAssessment();
    }, [currentAssessment]);

    return (<div className={b("wrapper")}>
        {
            !hasCurrentAssessment
                ? <>
                    {currentAssessment === "loading" && <Spin />}
                    {currentAssessment === "error" && 
                        <Alert title="Кончились задания 🤯" theme="danger" />}
                </>
                : <Flex direction="column" gap={4}>
                    <MarkupData assessment={currentAssessment} /> 
                    <MyMarkupType
                        isAdmin={false}
                        triggerRerender={() => { }} markupType={currentAssessment.markup_type}
                        onUpdateValue={setCurrentValue}
                    />
                    <Button view="action" onClick={async () => 
                        {
                            setCurrentAssessment("loading");
                            await assessmentUpdate(currentAssessment.assessment_id, valueToFields(currentValue));
                            setCurrentAssessment(null);

                            toaster.add({
                                name: "Отправлено",
                                content: "+0.03 ₽",
                                theme: "success",
                            })
                        }}
                    >
                        Отправить
                    </Button>
                </Flex>
        }
    </div>)
};
