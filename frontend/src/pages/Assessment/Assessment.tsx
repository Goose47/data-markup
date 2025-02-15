import { Button, Spin } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import { assessmentNext, assessmentUpdate } from "../../utils/requests";
import { AssessmentNext } from "../../utils/types";
import "./Assessment.scss";
import { useEffect, useState } from "react";
import { MyMarkupType } from "../../components/MyMarkupType/MyMarkupType";

const b = block("assessment");

const valueToFields = (assessment: AssessmentNext, value: string[][]) => {
    return [assessment, value]
    // const fieldTypeByValue = Object.fromEntries(assessment.markup_type.fields?.map((field) => [field.id, field.assessment_type_id]) ?? []);

    // return { fields: value.flatMap((group) => group.map((fieldValue) => { 
    //         return fieldTypeByValue[fieldValue] === 5 ? {
    //             markup_type_field_id: parseInt(fieldValue),
    //         } : {
    //         markup_type_field_id: parseInt(fieldValue),
    //     }}
    // ))}
};

export const Assessment = () => {
    const [currentAssessment, setCurrentAssessment] = useState<AssessmentNext | null>(null);
    const [currentValue, setCurrentValue] = useState<string[][]>([[]]);

    useEffect(() => {
        const fetchAssessment = async () => {
            if (currentAssessment !== null) {
                return;
            }

            const assessment: AssessmentNext = await assessmentNext();

            setCurrentAssessment(assessment);
        }

        fetchAssessment();
    }, [currentAssessment]);

    return (<div className={b("wrapper")}>
        {/* {
            currentAssessment === null 
                ? <Spin /> 
                : <MyMarkupType
                    isAdmin={false}
                    triggerRerender={() => { }} markupType={currentAssessment.markup_type}
                    onUpdateValue={setCurrentValue} 
                  />
        } */}
        <Button view="action" onClick={() => {
            assessmentUpdate(1, { fields: [] });
        }}>Отправить</Button>
    </div>)
};
