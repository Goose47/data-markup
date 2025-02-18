import { Modal, Table } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import "./UserStatPage.scss";
import { useContext, useEffect, useMemo, useState } from "react";
import { AsseessmentType } from "../../utils/types";
import { profileMe } from "../../utils/requests";
import { Assessment } from "../Assessment/Assessment";
import { LoginContext } from "../Login/LoginContext";

const b = block("user-stat-page");

export const UserStatPage = () => {
  const [assessments, setAssessments] = useState<AsseessmentType[]>([]);

  const [triggerRerender, setTriggerRerender] = useState(0);

  const handleTriggerRerender = () => {
    setTriggerRerender((s) => s + 1);
  };

  useEffect(() => {
    profileMe().then((data: { assessments: AsseessmentType[] }) => {
      setAssessments(data.assessments.filter((el) => el.hash !== null));
    });
  }, [triggerRerender]);

  const loginContext = useContext(LoginContext);

  const answers = useMemo(() => {
    return assessments.map((assessment) => {
      return assessment.fields
        .map((el) => el.markup_type_field)
        .map((x) => {
          if (x) {
            return `${x.label}: ${x.name}`;
          } else {
            return "-";
          }
        })
        .join(", ");
    });
  }, [assessments]);

  const assessmentsTableData = useMemo(() => {
    return assessments.map((el, index) => {
      return {
        answers: answers[index],
        is_correct_answer: el.is_correct ? "✅" : "❌",
        is_editable:
          loginContext.userRole === "admin" || el.is_editable ? "✅" : "❌",
        is_editable_raw: el.is_editable,
        markup_id: el.markup_id,
        assessment_id: el.id,
      };
    });
  }, [answers, assessments, loginContext.userRole]);

  const columns = useMemo(() => {
    return [
      { id: "markup_id", name: "ID разметки" },
      { id: "is_correct_answer", name: "Верный ответ" },
      { id: "is_editable", name: "Можно ли редактировать?" },
      { id: "answers", name: "Ответы" },
    ];
  }, []);

  const [showEdit, setShowEdit] = useState(false);
  const [markupId, setMarkupId] = useState(0);
  const [isEditable, setIsEditable] = useState(false);
  const [assessmentId, setAssessmentId] = useState(-1);

  return (
    <div className={b()}>
      <div className={b("block")}>
        <h1>Мои оценки</h1>
        <Table
          columns={columns}
          data={assessmentsTableData}
          onRowClick={(el) => {
            setMarkupId(el.markup_id ?? 0);
            setShowEdit(true);
            setIsEditable(el.is_editable_raw ?? false);
            setAssessmentId(el.assessment_id ?? -1);
          }}
        ></Table>
        {showEdit && (
          <Modal
            open={showEdit}
            onClose={() => setShowEdit(false)}
            contentClassName={b("popup")}
          >
            <Assessment
              markupId={markupId}
              assessmentId={assessmentId}
              isEditable={isEditable}
              triggerRerender={handleTriggerRerender}
            />
          </Modal>
        )}
      </div>
    </div>
  );
};
