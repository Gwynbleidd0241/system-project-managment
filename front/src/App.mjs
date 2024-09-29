import React, { useState, useEffect } from "react";
import './styles/App.css';
import TaskItem from "./components/TaskItem.jsx";
import MyButton from "./components/UI/button/MyButton.jsx";
import Modal from "./components/UI/modal/Modal.jsx";
import NewTaskModal from "./components/UI/newTask/NewTaskModal.jsx";
import RedTask from "./components/UI/redTask/RedTask.jsx";
import EnterModal from "./components/UI/EnterModal/EnterModal.jsx";
import axios from "axios";

function App() {
  const [tasks, setTasks] = useState([]);
  const maxTasks = 4;

  const [showAll, setShowAll] = useState(false);
  const [modal, setModal] = useState(false);
  const [isAuthorized, setIsAuthorized] = useState(false);
  const [newTaskModal, setNewTaskModal] = useState(false);
  const [redTaskModal, setRedTaskModal] = useState(false);
  const [currentTask, setCurrentTask] = useState(null);
  const [enterVisible, setEnterVisible] = useState(false);


  const handleShowAll = () => {
    setShowAll(!showAll);
  };

  const removeTask = (taskId) =>{
    console.log('Удаляем задачу с id:', taskId)
    setTasks([...tasks.filter(t => t.ID !== taskId)])
  }

  const onCreateTask = (task) => {
    setTasks([task, ...tasks]);
    setNewTaskModal(false);
  };

  const onEditTask = (task) => {
    setCurrentTask(task);
    setRedTaskModal(true);
  };

  const onSaveTask = (task) => {
    setTasks([...tasks.map(t => t.id === task.id ? task : t)]);
    setRedTaskModal(false);
  };


 return (
  <div className="App">
    <div className="head">
      <h1>ZOV</h1>
      {isAuthorized ? (
        <MyButton onClick={() => setIsAuthorized(false)}>Выйти</MyButton>
      ) : (
        <div>
          <MyButton onClick={() => setModal(true)}>Регистрация</MyButton>
          <MyButton onClick={() => setEnterVisible(true)}>Вход</MyButton>
        </div>
      )}
    </div>
    {modal && (<Modal visible={modal} setVisible={setModal} isAuthorized={isAuthorized} setIsAuthorized={setIsAuthorized} />)}
    {enterVisible && (<EnterModal visible={enterVisible} setVisible={setEnterVisible} isAuthorized={isAuthorized} setIsAuthorized={setIsAuthorized} setTasks={setTasks} />)}
    <NewTaskModal visible={newTaskModal} setVisible={setNewTaskModal} onCreateTask={onCreateTask} />
    <RedTask visible={redTaskModal} setVisible={setRedTaskModal} task={currentTask} onSaveTask={onSaveTask} onDeleteTask={removeTask} />
    {isAuthorized ? (
      <div className="tasks__container">
        <div className="Shower">
          {showAll ? (
            <MyButton onClick={handleShowAll}>Скрыть все</MyButton>
          ) : (
            <MyButton onClick={handleShowAll}>Показать все</MyButton>
           )}
        </div>
        <div className="tasks__fiels">
          {showAll ? (
            tasks.map((task, index) => {
              return <TaskItem 
                key={task.ID} 
                task={task} 
                number={index + 1} 
                onEdit={onEditTask} 
            />;
            })
          ) : (
            tasks.slice(0, maxTasks).map((task, index) => {
              return <TaskItem 
                key={task.ID} 
                task={task} 
                number={index + 1} 
                onEdit={onEditTask} 
            />;
            })
          )}
          <TaskItem
          task={{ ID: null, title: "Создать новую задачу", body: "" }}
          onClick={() => setNewTaskModal(true)}/>
        </div>
      </div>
      ) : (
        <p></p>
      )}
      <footer className="footer">
        <div className="footer__content">
          <p>&copy; 2024 SLON_ISLAND. Все права защищены<a href="https://contract.mos.ru/" target="_blank" className="footer__link">.</a></p>
        </div>
      </footer>
    </div>
  );
}

export default App;