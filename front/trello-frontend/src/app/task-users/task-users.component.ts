import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TaskService } from '../service/task.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-task-users',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './task-users.component.html',
  styleUrls: ['./task-users.component.css']
})
export class TaskUsersComponent implements OnInit {
  taskId!: string;
  users: any[] = []; // Zamenite `any` sa odgovarajućim modelom za korisnike

  constructor(
    private route: ActivatedRoute,
    private taskService: TaskService
  ) {}

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.taskId = params.get('taskId') || '';
      this.getUsersForTask();
    });
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
}
