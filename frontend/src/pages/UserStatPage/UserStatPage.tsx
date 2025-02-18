import { Loader, Pagination, PaginationProps, Table } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import "./UserStatPage.scss";
import { useContext, useEffect, useMemo, useState } from "react";
import { AsseessmentType } from "../../utils/types";
import { profileMe } from "../../utils/requests";
import { LoginContext } from "../Login/LoginContext";

const b = block("user-stat-page");

export const UserStatPage = () => {
  const [assessments, setAssessments] = useState<AsseessmentType[]>([]);
  const [state, setState] = useState({ page: 1, pageSize: 10 });
  const [totalPages, setTotalPages] = useState(1);
  const [correctHash, setCorrectHash] = useState("no set");

  const [triggerRerender, setTriggerRerender] = useState(0);

  const handleTriggerRerender = () => {
    setTriggerRerender((s) => s + 1);
  };

  const handleUpdate: PaginationProps["onUpdate"] = (page, pageSize) =>
    setState((prevState) => ({ ...prevState, page, pageSize }));

  const pagination = (
    <Pagination
      page={state.page}
      pageSize={state.pageSize}
      total={state.pageSize * totalPages}
      onUpdate={handleUpdate}
    />
  );

  useEffect(() => {
    profileMe().then((data: { assessments: AsseessmentType[] }) => {
      setAssessments(data.assessments);
    });
  }, []);

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
        is_editable: el.is_editable ? "✅" : "❌",
      };
    });
  }, [answers, assessments, correctHash]);

  const columns = useMemo(() => {
    return [
      { id: "is_correct_answer", name: "Верный ответ" },
      { id: "is_editable", name: "Можно ли редактировать?" },
      { id: "answers", name: "Ответы" },
    ];
  }, []);

  const loginContext = useContext(LoginContext);
  if (loginContext.loading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: 500,
        }}
      >
        <Loader></Loader>
      </div>
    );
  }
  if (loginContext.userRole !== "admin") {
    return (
      <div className={b()}>
        <h1>Отказано в доступе</h1>
      </div>
    );
  }

  return (
    <div className={b()}>
      <div className={b("block")}>
        <h1>Мои оценки</h1>
        <Table columns={columns} data={assessmentsTableData}></Table>
        {pagination}
      </div>
    </div>
  );
};
