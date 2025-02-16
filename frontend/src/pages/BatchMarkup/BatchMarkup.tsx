import { useNavigate, useParams } from "react-router";
import { block } from "../../utils/block";
import "./BatchMarkup.scss";
import { Loader, Pagination, PaginationProps, Table } from "@gravity-ui/uikit";
import { useEffect, useState } from "react";
import { getLinkedMarkupsToBatch } from "../../utils/requests";
import { BatchMarkupType } from "../../utils/types";

const b = block("batch-markup");

export const BatchMarkup = () => {
  const params = useParams();

  const navigate = useNavigate();

  const [state, setState] = useState({ page: 1, pageSize: 20 });

  const [data, setData] = useState<BatchMarkupType[]>([]);
  const [totalPages, setTotalPages] = useState(1);

  const handleUpdate: PaginationProps["onUpdate"] = (page, pageSize) =>
    setState((prevState) => ({ ...prevState, page, pageSize }));

  const batchId = parseInt(params.batchId ?? "0");

  useEffect(() => {
    if (isNaN(batchId)) {
      navigate("/");
      return;
    }

    getLinkedMarkupsToBatch(
      batchId,
      state.page,
      state.pageSize
    ).then((data: { data: BatchMarkupType[]; pages_total: number }) => {
      setData(data.data);
      setTotalPages(data.pages_total);
    });
  }, [batchId, navigate, params.batchId, state.page, state.pageSize]);

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

  return (
    <div className={b()}>
      <h1>АЗрара</h1>
      {parsedJson.length === 0 ? (
        <Loader />
      ) : (
        <Table
          columns={Object.keys(parsedJson[0]).map((k) => ({
            id: k,
            name: k,
          }))}
          data={data.map((_, index) =>
            Object.keys(parsedJson[index]).reduce(
              function (result, key) {
                result[key] = parsedJson[index][key].substring(0, 20) + "...";
                return result;
              },
              {} as any[string]
            )
          )}
        ></Table>
      )}

      {pagination}
    </div>
  );
};
