import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TaskService } from '../service/task.service';
import { CommonModule } from '@angular/common';
import { UserService } from '../service/user.service';
import { FormsModule } from '@angular/forms';
import { Task } from '../model/task.model';

@Component({
  selector: 'app-task-users',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './task-users.component.html',
  styleUrls: ['./task-users.component.css']
})
export class TaskUsersComponent implements OnInit {
  taskId!: string;
  task!: any;
  projectId!: string;
  users: any[] = []; 
  projectUsers: any[] = [];
  selectedUserId: string = '';
  showAddUserForm: boolean = false;

  constructor(
    private route: ActivatedRoute,
    private taskService: TaskService,
    private userService: UserService
  ) {}

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.taskId = params.get('taskId') || '';
      this.getTaskDetails();
    });
  }
  getTaskDetails(): void {
    if (this.taskId) {
      this.taskService.getTaskById(this.taskId).subscribe(
        (task) => {
          this.task = task;
          this.projectId = task.projectId; 
          this.getUsersForTask();
          this.getProjectUsers(); 
        },
        (error) => {
          console.error('Error fetching task details:', error);
        }
      );
    }
  }

  getUsersForTask(): void {
    if (this.taskId) {
      this.taskService.getUsersByTaskId(this.taskId).subscribe(
        (users) => {
          this.users = users;
        },
        (error) => {
          console.error('Error fetching users:', error);
        }
      );
    }
  }
  getProjectUsers(): void {
    if (this.projectId) {
      console.log("Project ID: ", this.projectId)
      this.userService.getUsersByProjectId(this.projectId).subscribe(
        (users) => {
          this.projectUsers = users.filter((user: { id: any; }) => !this.users.some(taskUser => taskUser.id === user.id));
          
        },
        (error) => {
          console.error('Error fetching project users:', error);
        }
      );
    }
  }
  addUserToTask(): void {
    if (this.selectedUserId) {
      this.taskService.addMemberToTask(this.taskId, this.selectedUserId).subscribe({
        next: () => {
          console.log(`User ${this.selectedUserId} added to task ${this.taskId}`);
          // Ažuriramo listu korisnika u zadatku nakon dodavanja
          this.getUsersForTask();
          this.getProjectUsers();
          this.showAddUserForm = false; // Zatvaramo formu nakon dodavanja
        },
        error: (error) => {
          console.error(`Error adding user to task:`, error);
        },
        complete: () => {
          console.log('User addition complete');
        }
      });
    }
  }
  removeUser(userId: string): void {
    if (this.taskId) {
      this.taskService.removeUserFromTask(this.taskId, userId).subscribe({
        next: () => {
          console.log(`User ${userId} removed from task ${this.taskId}`);
          this.users = this.users.filter(user => user.id !== userId); 
        },
        error: (error) => {
          console.error(`Error removing user from task:`, error);
        },
        complete: () => {
          console.log('User removal complete');
        }
      });
    }
  }
  isUserManager(): boolean{
    var role = localStorage.getItem('userRole')
    if(role === 'manager'){
      return true
      }
      return false
  }
}
