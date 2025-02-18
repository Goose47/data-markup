import { useContext, useEffect, useState } from "react";
import { block } from "../../utils/block";
import "./MyMarkupTypes.scss";
import { MarkupType } from "../../utils/types";
import { getAvailableMarkupTypes } from "../../utils/requests";
import { MyMarkupType } from "../../components/MyMarkupType/MyMarkupType";
import { CircleInfoFill } from "@gravity-ui/icons";
import { LoginContext } from "../Login/LoginContext";
import { Loader } from "@gravity-ui/uikit";

const b = block("my-markup-types");

export const MyMarkupTypes = () => {
  const [markupTypes, setMarkupTypes] = useState<MarkupType[]>([]);

  const [rerenderState, setRerenderState] = useState(1);
  useEffect(() => {
    getAvailableMarkupTypes().then((value: MarkupType[]) => {
      setMarkupTypes(value);
    });
  }, [rerenderState]);

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
      <h1>Существующие типы разметок</h1>
      <p>
        <CircleInfoFill></CircleInfoFill> Здесь отображены те типы, которые
        можно использовать при создании проектов (batch'ей). На данной странице
        вы можете проверить, как будет выглядеть со стороны ассессора варианты
        ответа на определенном задании, куда вы добавите этот тип.
      </p>

      <div className={b("list")}>
        {markupTypes.map((markupType) => (
          <MyMarkupType
            isAdmin={true}
            markupType={markupType}
            triggerRerender={() => setRerenderState((v) => v + 1)}
          />
        ))}
      </div>
    </div>
  );
};
