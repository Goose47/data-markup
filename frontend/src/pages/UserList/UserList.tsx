import { useContext, useEffect, useMemo, useState } from "react";
import { block } from "../../utils/block";
import "./UserList.scss";
import { UserListType } from "../../utils/types";
import { getAllProfiles } from "../../utils/requests";
import { Loader, Table } from "@gravity-ui/uikit";
import { LoginContext } from "../Login/LoginContext";

const b = block("user-list");

export const UserList = () => {
  const [users, setUsers] = useState<UserListType[]>([]);

  useEffect(() => {
    getAllProfiles().then((data: UserListType[]) => {
      if (!data.length) return;
      setUsers(
        data?.map((el) => {
          return {
            ...el,
            created_at: new Date(el.created_at).toLocaleString("ru-RU"),
            precision:
              el.assessment_count + el.assessment_count2 !== 0
                ? (
                    (100 *
                      (el.correct_assessment_count +
                        el.correct_assessment_count2)) /
                    (el.assessment_count + el.assessment_count2)
                  ).toFixed(2) + "%"
                : "-",
          };
        })
      );
    });
  }, []);

  const columns = useMemo(() => {
    return [
      { id: "email", name: "E-mail" },
      { id: "precision", name: "Точность ответов" },
      {
        id: "assessment_count",
        name: "Ответов по поиску",
      },
      { id: "correct_assessment_count", name: "Из них верных" },
      { id: "assessment_count2", name: "Ответов по сравнению" },
      { id: "correct_assessment_count2", name: "Из них верных" },
      { id: "created_at", name: "Зарегистрирован" },
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
      <h1>Общий список ассессоров и администраторов</h1>
      <Table columns={columns} data={users}></Table>
    </div>
  );
};
