import React, { useState, useEffect, useCallback } from "react";
import axios from 'axios';
import cl from './EnterModal.module.css'

const EnterModal = ({ visible, setVisible, isAuthorized, setIsAuthorized, setTasks }) => {
    const rootClasses = [cl.modal]
  if (visible){
    rootClasses.push(cl.active);
  }
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState(null);
  const [token, setToken] = useState(null);

  const handleClose = useCallback(() => {
    setVisible(false);
  }, [setVisible]);

  const handleSubmit = (e) => {
    e.preventDefault();
    axios.post('http://localhost:8080/api/login', {
        email: login,
        password: password,
    })
    .then(response => {
        const token = response.data.token;
        setToken(token);
        localStorage.setItem('token', token);
        localStorage.setItem("login", login);
        setIsAuthorized(true);
        setVisible(false);

        // Второй запрос после успешной авторизации
        return axios.get('http://localhost:8080/api/task', {
          params: {
            email: login
          },
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
    })
    .then(response => {
      const tasks = response.data;
      console.log(tasks) // Обеспечьте, чтобы tasks был массивом
      setTasks(tasks);
        localStorage.setItem('login', login);
    })
    .catch(error => {
        setError('Ошибка входа. Сервер не отвечает.');
        // setTimeout(() => {
        //     handleClose(); // закрываем модальное окно после успешной регистрации
        // }, 2000);
    });
};

  useEffect(() => {
    if (token) {
      console.log(login)
      axios.get('http://localhost:8080/api/task', {
        params: {
          email: login
        },
        headers: {
          Authorization: `Bearer ${token}`
        }
      })
      .then(response => {
        const tasks = response.data.tasks;
        setTasks(tasks);
      })
      .catch(error => {
        console.error(error);
      });
    }
  }, [token, login, setTasks]);

  if (!visible) return null;

  return (
    <div className={rootClasses.join(" ")} onClick={handleClose}>
      <div className={cl.modalContent} onClick={(e) => e.stopPropagation()}>
        <form onSubmit={handleSubmit}>
          <label>
            Логин:
            <input type="text" value={login} onChange={(e) => setLogin(e.target.value)} />
          </label>
          <label>
            Пароль:
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
          </label>
          {error && <div style={{ color: 'red' }}>{error}</div>}
          <button type="submit">Войти</button>
        </form>
      </div>
    </div>
  );
};

export default EnterModal;