#include <fcntl.h>
#include <semaphore.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>
#include "operations.h"

// Compile with flags: -lrt -lpthreads

/**
 * Friendly names for supported worker process operations.
 */
char *op_names[] = {
        "add", "sub", "mul", "div", "quit"
};

bool create_shared_object(shared_memory_t *shm, const char *share_name) {
    // Remove any previous instance of the shared memory object, if it exists.
    shm_unlink(share_name);

    // Assign share name to shm->name.
    shm->name = share_name;

    // Create the shared memory object, allowing read-write access, and saving the
    // resulting file descriptor in shm->fd. If creation failed, ensure
    // that shm->data is NULL and return false.
    if ((shm->fd = shm_open(share_name, O_RDWR | O_CREAT)) < 0) {
        shm->data = NULL;
        return false;
    }
    // Set the capacity of the shared memory object via ftruncate. If the
    // operation fails, ensure that shm->data is NULL and return false.
    if (ftruncate(shm->fd, sizeof(shared_data_t)) < 0) {
        shm->data = NULL;
        return false;
    }

    // Otherwise, attempt to map the shared memory via mmap, and save the address
    // in shm->data. If mapping fails, return false.
    shm->data = (shared_data_t*)malloc(sizeof(shared_data_t));
    if (mmap(shm->data, sizeof(shared_data_t), PROT_READ | PROT_WRITE, MAP_SHARED, shm->fd, 0) == MAP_FAILED) {
        return false;
    }
    // Do not alter the following semaphore initialisation code.
    sem_init(&shm->data->controller_semaphore, 1, 0);
    sem_init(&shm->data->worker_semaphore, 1, 0);

    // If we reach this point we should return true.
    return true;
}

void destroy_shared_object(shared_memory_t *shm) {
    // Remove the shared memory object.
    munmap(shm->data, sizeof(shared_data_t));
    shm_unlink(shm->name);
    shm->fd = -1;
    shm->data = NULL;
}

double request_work(shared_memory_t *shm, operation_t op, double lhs, double rhs) {
    // Copy the supplied values of op, lhs and rhs into the corresponding fields
    // of the shared data object.
    shm->data->operation = op;
    shm->data->lhs = lhs;
    shm->data->rhs = rhs;

    // Do not alter the following semaphore code. It sends the request to the
    // worker, and waits for the response in a reliable manner.
    sem_post(&shm->data->controller_semaphore);

    sem_wait(&shm->data->worker_semaphore);

    // Modify the following line to make the function return the result computed
    // by the worker process. This will be stored in the result field of the
    // shared data object.

    return shm->data->result;
}

bool get_shared_object(shared_memory_t *shm, const char *share_name) {
    // Get a file descriptor connected to shared memory object and save in
    // shm->fd. If the operation fails, ensure that shm->data is
    // NULL and return false.
    if ((shm->fd = shm_open(share_name, O_RDWR)) < 0) {
        shm->data = NULL;
        return false;
    }

    // Otherwise, attempt to map the shared memory via mmap, and save the address
    // in shm->data. If mapping fails, return false.
    shm->data = (shared_data_t*)malloc(sizeof(shared_data_t));
    if (mmap(shm->data, sizeof(shared_data_t), PROT_READ | PROT_WRITE, MAP_SHARED, shm->fd, 0) == MAP_FAILED) {
        return false;
    }

    // Modify the remaining stub only if necessary.
    return true;
}


bool do_work(shared_memory_t *shm) {
    bool retVal = true;

    // Do not alter the following instruction, which waits for work
    sem_wait(&shm->data->controller_semaphore);

    // Update the value of local variable retVal and/or shm->data->result
    // as required.

    retVal = (shm->data->operation != op_quit);
    if (retVal) {
        double result = 0;
        switch (shm->data->operation) {
            case op_add:
                result = shm->data->lhs + shm->data->rhs;
                break;
            case op_sub:
                result = shm->data->lhs - shm->data->rhs;
                break;
            case op_mul:
                result = shm->data->lhs * shm->data->rhs;
                break;
            case op_div:
                result = shm->data->lhs / shm->data->rhs;
                break;
            default:
                break;
        }
        shm->data->result=result;
    }

    // Do not alter the following instruction which send the result back to the
    // controller.
    sem_post(&shm->data->worker_semaphore);

    // If retval is false, the memory needs to be unmapped, but that must be
    // done _after_ posting the semaphore. Un-map the shared data, and assign
    // values to shm->data and shm-fd as noted above.
    // INSERT IMPLEMENTATION HERE
    if (!retVal) {
        munmap(shm->data, sizeof(shared_data_t));
        shm->fd = -1;
        shm->data = NULL;
        return retVal;
    }

    // Keep this line to return the result.
    return retVal;
}

double next_rand() {
    return 100.0 * rand() / RAND_MAX;
}

operation_t next_op() {
    return (operation_t) (rand() % op_quit);
}

#define SHARE_NAME "/xyzzy_123"

void controller_main() {
    srand(42);
    printf("Controller starting.\n");

    shared_memory_t shm;

    if (create_shared_object(&shm, SHARE_NAME)) {
        for (int i = 0; i < 20; i++) {
            operation_t op = next_op();
            double lhs = next_rand();
            double rhs = next_rand();
            double result = request_work(&shm, op, lhs, rhs);
            printf("%s(%0.2f, %0.2f) = %0.2f\n", op_names[op], lhs, rhs, result);
        }

        request_work(&shm, op_quit, 0, 0);
        printf("Controller finished.\n");

        destroy_shared_object(&shm);
    } else {
        printf("Shared memory creation failed.\n");
    }
}

void worker_main() {
    printf("Worker starting.\n");

    shared_memory_t shm;

    if (get_shared_object(&shm, SHARE_NAME)) {
        while (do_work(&shm)) {}

        printf("Worker has been told to quit.\n");
    } else {
        printf("Shared memory connection failed.\n");
    }
}

int main() {
    pid_t childPid = -1;
    // Invoke the fork function to spawn the worker process, and save the result
    // as childPid.
    childPid = fork();

    if (childPid < 0) { /* error occurred */
        fprintf(stderr, "Fork failed\n");
        return 1;
    } else if (childPid == 0) {
        // Sleep 1 second to give the controller time to create the shared memory
        // object, then invoke worker_main.
        sleep(1);
        worker_main();
    } else { /* parent process */
        // Invoke controller_main.
        controller_main();
    }

    return 0;
}