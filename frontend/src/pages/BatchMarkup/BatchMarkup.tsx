import { useNavigate, useParams } from "react-router";
import { block } from "../../utils/block";
import "./BatchMarkup.scss";
import {
  Button,
  Loader,
  Pagination,
  PaginationProps,
  Progress,
  Table,
} from "@gravity-ui/uikit";
import { useEffect, useState } from "react";
import {
  batchFind,
  downloadBatchResult,
  getLinkedMarkupsToBatch,
} from "../../utils/requests";
import { BatchMarkupType } from "../../utils/types";
import { BatchHistory } from "../../components/BatchHistory/BatchHistory";

const b = block("batch-markup");

export const BatchMarkup = () => {
  const params = useParams();

  const [state, setState] = useState({ page: 1, pageSize: 5 });

  const [batch, setBatch] = useState<BatchMarkupType>();
  const [data, setData] = useState<BatchMarkupType[]>([]);
  const [totalPages, setTotalPages] = useState(1);

  const handleUpdate: PaginationProps["onUpdate"] = (page, pageSize) =>
    setState((prevState) => ({ ...prevState, page, pageSize }));

  useEffect(() => {
    getLinkedMarkupsToBatch(
      parseInt(params.batchId ?? "0"),
      state.page,
      state.pageSize
    ).then((data: { data: BatchMarkupType[]; pages_total: number }) => {
      setData(data.data);
      setTotalPages(data.pages_total);
    });
    batchFind(parseInt(params.batchId ?? "0")).then((data: BatchMarkupType) => {
      setBatch(data);
    });
  }, [params.batchId, state.page, state.pageSize]);

  const pagination = (
    <Pagination
      page={state.page}
      pageSize={state.pageSize}
      total={state.pageSize * totalPages}
      onUpdate={handleUpdate}
    />
  );

  const parsedJson = data.map((markup: BatchMarkupType) =>
    JSON.parse(markup.data)
  );

  const navigate = useNavigate();

  return (
    <div className={b()}>
      <h1>Прогресс по проекту</h1>
      {!batch ? (
        <Loader></Loader>
      ) : (
        <>
          <Progress
            text={
              (
                (batch.processed_markup_count / batch.markup_count) *
                100
              ).toFixed(2) + "%"
            }
            value={(batch.processed_markup_count / batch.markup_count) * 100}
            colorStops={[
              { theme: "danger", stop: 20 },
              { theme: "warning", stop: 50 },
              { theme: "success", stop: 100 },
            ]}
          />
          {batch.assessment_count ? (
            <h1>
              Точность ассессоров на проекте:{" "}
              {(
                (100 * batch.correct_assessment_count) /
                batch.assessment_count
              ).toFixed(2)}
              %
            </h1>
          ) : (
            <></>
          )}

          <p>
            Количество оценок: {batch.assessment_count}
            <br />
            Количество корректных оценок: {batch.correct_assessment_count}
            <br />
            Количество раметок: {batch.markup_count}
            <br />
            Количество завершенных разметок: {batch.processed_markup_count}
            <br></br>
            <br></br>{" "}
            <Button
              view="action"
              onClick={() => {
                downloadBatchResult(parseInt(params.batchId ?? "0"));
              }}
            >
              Скачать результат
            </Button>
          </p>
        </>
      )}

      <h1>Данные, привязанные к вашему проекту</h1>

      {parsedJson.length === 0 ? (
        <Loader />
      ) : (
        <Table
          onRowClick={(row) => {
            navigate("/batch/markup/" + row.id);
          }}
          columns={[
            { id: "id", name: "ID" },
            { id: "number_of_answers", name: "Количество ответов" },
            { id: "correct_answer", name: "Есть верный ответ" },
            ...Object.keys(parsedJson[0]).map((k) => ({
              id: k,
              name: k,
            })),
          ]}
          data={[
            ...data.map((_, index) => ({
              id: data[index].id,
              correct_answer:
                data[index].correct_assessment_hash !== null ? "да" : "нет",
              number_of_answers: data[index].assessments?.length ?? 0,
              ...Object.keys(parsedJson[index]).reduce(
                function (result, key) {
                  result[key] =
                    parsedJson[index][key].substring(0, 30) +
                    (parsedJson[index][key] > 30 ? "..." : "");
                  return result;
                },
                {} as any[string]
              ),
            })),
          ]}
        ></Table>
      )}
      {pagination}

      <h1>История изменения типа разметки</h1>
      <BatchHistory batchId={parseInt(params.batchId ?? "0")}></BatchHistory>
    </div>
  );
};
