import { useParams } from "react-router";
import { block } from "../../utils/block";
import "./MarkupAssessments.scss";
import {
  Button,
  Loader,
  Pagination,
  PaginationProps,
  Table,
} from "@gravity-ui/uikit";
import { useContext, useEffect, useMemo, useState } from "react";
import { Assessment } from "../Assessment/Assessment";
import { LoginContext } from "../Login/LoginContext";
import { AsseessmentType } from "../../utils/types";
import {
  assessmentIndex,
  getBatchMarkupData,
  makeHoneypot,
} from "../../utils/requests";
import { ButtonWithConfirm } from "../../components/ButtonWithConfirm/ButtonWithConfirm";
import { CircleInfoFill } from "@gravity-ui/icons";

const b = block("markup-assessments");

export const MarkupAssessments = () => {
  const params = useParams();

  const markupId = params.markupId ?? "0";

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
    assessmentIndex(parseInt(markupId), state.page, state.pageSize).then(
      (data: { data: AsseessmentType[]; pages_total: number }) => {
        setAssessments(data.data);
        setTotalPages(data.pages_total);
      }
    );
    getBatchMarkupData(parseInt(markupId)).then((data) => {
      setCorrectHash(data.correct_assessment_hash);
    });
  }, [markupId, state.page, state.pageSize, triggerRerender]);

  const answers = useMemo(() => {
    return assessments.map((assessment) => {
      return assessment.fields
        .map((el) => el.markup_type_field_id)
        .map((id) => {
          const x = assessment.markup_type.fields.find((el) => el.id === id);
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
        user_email: el.user.email,
        answers: answers[index],
        is_correct_answer: el.hash === correctHash ? "✅" : "❌",
        is_admin: el.is_prior ? "✅" : "❌",
      };
    });
  }, [answers, assessments, correctHash]);

  const columns = useMemo(() => {
    return [
      { id: "user_email", name: "Ассессор" },
      { id: "is_correct_answer", name: "Верный ответ" },
      { id: "is_admin", name: "Эталонная оценка" },
      { id: "answers", name: "Ответы" },
    ];
  }, []);

  const handleMakeHoneypot = () => {
    makeHoneypot(parseInt(markupId ?? "0")).then(() => {
      localStorage.setItem("honeypot" + markupId, "true");
    });
  };

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
        <ButtonWithConfirm
          handleSubmit={handleMakeHoneypot}
          confirmText="Вы действительно хотите добавить этот маркап всем ассессорам?"
          disabled={
            !Boolean(assessments.filter((el) => el.is_prior).length) ||
            Boolean(localStorage.getItem("honeypot" + markupId))
          }
        >
          <Button
            view="action"
            disabled={
              !Boolean(assessments.filter((el) => el.is_prior).length) ||
              Boolean(localStorage.getItem("honeypot" + markupId))
            }
          >
            Сделать ханипотом
          </Button>
          <p>
            <CircleInfoFill /> Чтобы сделать разметку ханипотом, необходимо,
            чтобы у него была оценка от администратора
            <br />
          </p>
        </ButtonWithConfirm>
        <h1>Таблица оценок</h1>
        <Table columns={columns} data={assessmentsTableData}></Table>
        {pagination}
      </div>
      <div className={b("block")}>
        <h1>Оценить от имени администратора</h1>
        <Assessment
          markupId={parseInt(markupId)}
          isAdmin={true}
          triggerRerender={handleTriggerRerender}
        />
      </div>
    </div>
  );
};
